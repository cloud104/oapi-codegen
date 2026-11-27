#!/usr/bin/env bash

set -xe

#find examples/ -type f -name '*.go' ! -name '*_test.go' -delete
go test ./cmd/oapi-codegen/ -write-results
#kts go clean examples/
#kts go format --disable-gofumpt --disable-imports-reviser --disable-wsl --enable-callorder examples/oauth2-client/client.gen.go
#kts go format --disable-gofumpt --disable-imports-reviser --disable-wsl --enable-goreorder examples/oauth2-client/client.gen.go
#kts go format examples/
