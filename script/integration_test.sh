#!/usr/bin/env bash

cd $(dirname $0)/..

set -x

gcloud beta emulators datastore start --host-port="localhost:1234" --no-store-on-disk &
dockerize --timeout 10s -wait http://localhost:1234
export DATASTORE_EMULATOR_HOST=localhost:1234
go test -v ./test/...
kill $!
