#!/usr/bin/env bash

set -eu

repo="$( cd $(dirname $0) && pwd )/../.."

export GOPATH="$repo/../../../.."
export PATH=${GOPATH}/bin:${PATH}

cd ${repo}

echo -e "\n Vetting packages for potential issues..."
./bin/govet

./bin/install-ginkgo

./bin/test-unit
