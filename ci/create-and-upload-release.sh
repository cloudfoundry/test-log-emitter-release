#!/usr/bin/env bash
set -euo pipefail

export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=$PASSWORD
export BOSH_CA_CERT=$CA_CERT

push test-log-emitter-release
  bosh create-release --force
  bosh clean-up
  bosh -n -e $TARGET upload-release
popd
