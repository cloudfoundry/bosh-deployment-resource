#!/usr/bin/env bash

set -ue

my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
    source utils.sh
    set_env
    gcloud_login
popd > /dev/null

pushd ${release_dir}
    trap clean_gcs EXIT
    GOOGLE_SERVICE_ACCOUNT=$google_json_key_data make test-int
popd