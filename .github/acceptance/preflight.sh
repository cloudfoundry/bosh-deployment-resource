#!/usr/bin/env bash
set -euo pipefail

if [[ "$(uname -s)" != Linux || "$(uname -m)" != x86_64 ]]; then
  printf '%s\n' 'This acceptance test requires a linux/amd64 host with a local Docker daemon.' >&2
  exit 1
fi

for tool in docker curl jq sha256sum; do
  command -v "$tool" >/dev/null
done
docker info >/dev/null

docker_root=$(docker info --format '{{.DockerRootDir}}')
available_kib=$(df -Pk "$docker_root" | awk 'NR == 2 {print $4}')
if (( available_kib < 12 * 1024 * 1024 )); then
  printf '%s\n' 'At least 12 GiB of free Docker disk space is required.' >&2
  exit 1
fi

# instant-bosh uses fixed resource names and ports. Never touch an existing environment.
if [[ -n "$(docker ps -aq --filter 'name=^/instant-bosh$')" ]] \
  || docker network inspect instant-bosh >/dev/null 2>&1 \
  || docker volume inspect instant-bosh-store >/dev/null 2>&1 \
  || docker volume inspect instant-bosh-data >/dev/null 2>&1; then
  printf '%s\n' 'An instant-bosh container, network, or volume already exists. Use a clean host.' >&2
  exit 1
fi
