#!/usr/bin/env bash

set -e

export BUILT_BINARIES_DIR=$PWD/built-binaries
export GOPATH="${PWD}/gopath"
export PATH="${GOPATH}/bin:${PATH}"
go get github.com/onsi/ginkgo/ginkgo

cd gopath/src/github.com/cloudfoundry/bosh-deployment-resource

go build -o "$BUILT_BINARIES_DIR/check" cmd/check
go build -o "$BUILT_BINARIES_DIR/in" cmd/in
go build -o "$BUILT_BINARIES_DIR/out" cmd/out
