#!/bin/bash

fly -t director set-pipeline -p davcli -c ci/pipeline.yml \
  -l <(lpass show -G "davcli concourse secrets" --notes) \
  -l <(lpass show --notes "pivotal-tracker-resource-keys")
