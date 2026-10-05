#!/usr/bin/env bash

set -euo pipefail

rm -rf go.mod go.sum
go mod init github.com/cloud104/oapi-codegen
go get \
	github.com/Masterminds/sprig/v3@latest \
	github.com/getkin/kin-openapi@latest \
	github.com/oapi-codegen/oapi-codegen/v2@latest \
	go.yaml.in/yaml/v3@latest \
	github.com/oapi-codegen/runtime@latest \
	golang.org/x/crypto@latest \
	golang.org/x/mod@latest \
	golang.org/x/sync@latest \
	golang.org/x/text@latest \
	golang.org/x/tools@latest \
  && go mod tidy
