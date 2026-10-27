#!/usr/bin/env bash
set -euo pipefail

repo_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
work_dir=$(mktemp -d)
started=false
export RESOURCE_IMAGE="bosh-deployment-resource:acceptance"
export IBOSH="$work_dir/ibosh"
export BOSH_BIN="$work_dir/bosh"
# Pin to the image that was verified together with the v0.3 ibosh release.
export INSTANT_BOSH_IMAGE="ghcr.io/rkoster/instant-bosh@sha256:369a5250348f9a9b0882e284b8d29cdfa234940a5672f20bc00c6b692d792222"

cleanup() {
  local result=$?
  trap - EXIT
  set +e
  if [[ "$started" == true ]]; then
    if [[ "$result" != 0 ]]; then
      "$IBOSH" docker logs --tail 120
    fi
    "$IBOSH" docker destroy --force
  fi
  rm -rf "$work_dir"
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

bash "$repo_dir/ci/acceptance/preflight.sh"

docker build --platform linux/amd64 \
  --build-arg "go_version=$(awk '/^go / {print $2}' "$repo_dir/go.mod")" \
  -f "$repo_dir/ci/acceptance/Dockerfile" -t "$RESOURCE_IMAGE" "$repo_dir"

curl --fail --location --retry 3 \
  https://github.com/rkoster/instant-bosh/releases/download/v0.3/instant-bosh_Linux_x86_64.tar.gz \
  -o "$work_dir/ibosh.tar.gz"
printf '%s  %s\n' 34a4ac291a07a8a285886523477a779cfcd11c86b49e8e2b0c2a28afcdd70a57 \
  "$work_dir/ibosh.tar.gz" | sha256sum --check
tar -xzf "$work_dir/ibosh.tar.gz" -C "$work_dir" ibosh

curl --fail --location --retry 3 \
  https://github.com/cloudfoundry/bosh-cli/releases/download/v7.11.0/bosh-cli-7.11.0-linux-amd64 \
  -o "$BOSH_BIN"
printf '%s  %s\n' 3ad84c3967672070b7f70f91f8da76abcbb604d588787b4f98457338cc9fc8d2 \
  "$BOSH_BIN" | sha256sum --check
chmod +x "$BOSH_BIN"

started=true
timeout --kill-after=30 1500 "$IBOSH" docker start --image "$INSTANT_BOSH_IMAGE"
timeout --kill-after=30 1200 bash "$repo_dir/ci/acceptance/test-resource.sh"
