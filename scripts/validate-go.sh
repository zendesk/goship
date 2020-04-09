#!/usr/bin/env bash

set -euo pipefail

echo "==> Running linters <=="
golangci-lint run -E gofmt -E stylecheck
exit $?
