#!/usr/bin/env bash
set -euo pipefail

docker info >/dev/null
docker run --rm --platform linux/amd64 --entrypoint bash bosh/docker-cpi:main -ec '
  command -v bosh
  command -v docker
  command -v jq
  test -f /usr/local/bosh-deployment/bosh.yml
  test -f /usr/local/bosh-deployment/docker/cpi.yml
  bosh --version
  docker --version
'