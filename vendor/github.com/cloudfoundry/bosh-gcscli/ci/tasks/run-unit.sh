#!/usr/bin/env bash

set -ue

my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
    source utils.sh
    set_env
popd > /dev/null

pushd ${release_dir}
    make test-unit
popd > /dev/null