#!/usr/bin/env bash
set -euo pipefail

: "${RESOURCE_IMAGE:?Set the local candidate image name}"
: "${IBOSH:?Set the ibosh binary path}"
: "${BOSH_BIN:?Set the bosh CLI path}"
repo_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT
mkdir -p "$test_dir/sources" "$test_dir/download"

# print-env exports Director credentials and a temporary jumpbox key path.
eval "$("$IBOSH" docker print-env)"
test -n "${BOSH_ALL_PROXY:-}"
proxy_key_path=${BOSH_ALL_PROXY#*private-key=}
test -f "$proxy_key_path"

cp "$repo_dir/ci/acceptance/deployment.yml" "$test_dir/sources/deployment.yml"
curl --fail --location --retry 3 \
  https://github.com/vmware-archive/dummy-boshrelease/releases/download/v2/dummy-2.tgz \
  -o "$test_dir/sources/release.tgz"
printf '%s  %s\n' d16c7d4a60306e890743dc3662a209b2c6c249d0e28e2639184c3bb6e79689fe \
  "$test_dir/sources/release.tgz" | sha256sum --check
jq -n '[{type: "replace", path: "/instance_groups/name=test/properties/marker",
  value: "((test_marker))-ops"}]' > "$test_dir/sources/ops.json"

jq -n \
  --arg target "$BOSH_ENVIRONMENT" --arg client "$BOSH_CLIENT" \
  --arg secret "$BOSH_CLIENT_SECRET" --arg ca "$BOSH_CA_CERT" \
  --rawfile key "$proxy_key_path" \
  '{deployment: "resource-acceptance", target: $target, client: $client,
    client_secret: $secret, ca_cert: $ca, jumpbox_url: "127.0.0.1:2222",
    jumpbox_username: "jumpbox", jumpbox_ssh_key: $key}' > "$test_dir/source.json"
chmod 600 "$test_dir/source.json"

resource() {
  docker run --rm -i --platform linux/amd64 --network host \
    -v "$test_dir:/work" "$RESOURCE_IMAGE" "/opt/resource/$1" "${@:2}"
}

jq -n --slurpfile source "$test_dir/source.json" '{source: $source[0], params: {
  manifest: "deployment.yml", releases: ["release.tgz"],
  ops_files: ["ops.json"], vars: {test_marker: "resource-acceptance"}
}}' | resource out /work/sources > "$test_dir/out-response.json"
jq -e '.version.manifest_sha1 | test("^[a-f0-9]{40}$")' "$test_dir/out-response.json"

"$BOSH_BIN" -d resource-acceptance --json instances > "$test_dir/instances.json"
jq -e '[.Tables[].Rows[]] | length == 1 and all(.[]; .process_state == "running")' \
  "$test_dir/instances.json"

jq -n --slurpfile source "$test_dir/source.json" --slurpfile out "$test_dir/out-response.json" \
  '{source: $source[0], version: $out[0].version}' > "$test_dir/version-request.json"
resource check < "$test_dir/version-request.json" > "$test_dir/check-response.json"
jq -e '. == []' "$test_dir/check-response.json"
resource in /work/download < "$test_dir/version-request.json" > "$test_dir/in-response.json"
jq -e --slurpfile out "$test_dir/out-response.json" '.version == $out[0].version' \
  "$test_dir/in-response.json"
test -s "$test_dir/download/manifest.yml"
test "$(cat "$test_dir/download/target")" = "$BOSH_ENVIRONMENT"
"$BOSH_BIN" interpolate "$test_dir/download/manifest.yml" \
  --path=/instance_groups/name=test/properties/marker | grep -Fx resource-acceptance-ops

jq -n --slurpfile source "$test_dir/source.json" \
  '{source: $source[0], params: {delete: {enabled: true}}}' \
  | resource out /work/sources > "$test_dir/delete-response.json"
"$BOSH_BIN" --json deployments > "$test_dir/deployments.json"
jq -e '[.Tables[].Rows[] | select(.name == "resource-acceptance")] | length == 0' \
  "$test_dir/deployments.json"
printf '%s\n' 'PASS: candidate image deploy, instance health, check, in, interpolation, and deletion.'
