#!/usr/bin/env bash

set -euo pipefail

exit_code=0

if ! hash gometalinter.v1 2>/dev/null ; then
  go get -u gopkg.in/alecthomas/gometalinter.v1
  gometalinter.v1 --install
fi

echo
echo "==> Running static validations <=="
# Run linters that should return errors
gometalinter.v1 \
  --disable-all \
  --enable deadcode \
  --severity deadcode:error \
  --enable gofmt \
  --enable ineffassign \
  --enable misspell \
  --enable vet \
  --tests \
  --vendor \
  --deadline 60s \
  ./... || exit_code=1

echo
echo "==> Running linters <=="
# Run linters that should return warnings
gometalinter.v1 \
  --disable-all \
  --enable golint \
  --vendor \
  --skip proto \
  --deadline 60s \
  ./... || :

exit $exit_code
