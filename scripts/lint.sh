#!/usr/bin/env bash

set -euo pipefail

go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run \
  --enable-only wrapcheck \
  --tests=false \
  --issues-exit-code 0 \
  --max-issues-per-linter 0 \
  --max-same-issues 0 \
  ./cmd/oapi-codegen/generated/...
