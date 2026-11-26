#!/usr/bin/env bash

set -xe

#find examples/ -type f -name '*.go' ! -name '*_test.go' -delete
go test ./cmd/oapi-codegen/ -write-results
#kts go clean examples/
#kts go format --disable-gofumpt --disable-imports-reviser --disable-wsl --enable-callorder examples/
#kts go format --disable-gofumpt --disable-imports-reviser --disable-wsl --enable-goreorder examples/
kts go format examples/
