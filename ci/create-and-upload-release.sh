#!/usr/bin/env bash
set -euo pipefail

export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=$PASSWORD
export BOSH_CA_CERT=$CA_CERT

pushd test-log-emitter-release
  bosh -n -e $TARGET create-release --force
  bosh -n -e $TARGET clean-up
  bosh -n -e $TARGET upload-release
popd
