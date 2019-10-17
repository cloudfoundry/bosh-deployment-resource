#!/usr/bin/env bash

check_param() {
  local name=$1
  local value=$(eval echo '$'$name)
  if [ "$value" == 'replace-me' ]; then
    echo "environment variable $name must be set"
    exit 1
  fi
}

print_git_state() {
  echo "--> last commit..."
  TERM=xterm-256color git log -1
  echo "---"
  echo "--> local changes (e.g., from 'fly execute')..."
  TERM=xterm-256color git status --verbose
  echo "---"
}

declare -a on_exit_items
on_exit_items=()

function on_exit {
  echo "Running ${#on_exit_items[@]} on_exit items..."
  for i in "${on_exit_items[@]}"
  do
    for try in $(seq 0 9); do
      sleep $try
      echo "Running cleanup command $i (try: ${try})"
        eval $i || continue
      break
    done
  done
}

function add_on_exit {
  local n=${#on_exit_items[@]}
  on_exit_items=("${on_exit_items[@]}" "$*")
  if [[ $n -eq 0 ]]; then
    trap on_exit EXIT
  fi
}

function clean_gcs {
    pushd ${release_dir}
        make clean-gcs
    popd
}

function set_env {
    my_dir=$(dirname "$(readlink -f "$0")")
    export release_dir="$( cd ${my_dir} && cd ../.. && pwd )"
    export workspace_dir="$( cd ${release_dir} && cd ../../../.. && pwd )"

    export GOPATH=${workspace_dir}
    export PATH=${GOPATH}/bin:${PATH}
}

function gcloud_login {
    check_param 'google_project'
    check_param 'google_json_key_data'

    keyfile=$(mktemp)
    gcloud config set project ${google_project}
    echo ${google_json_key_data} > ${keyfile}
    gcloud auth activate-service-account --key-file=${keyfile}
}