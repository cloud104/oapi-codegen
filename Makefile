.DEFAULT_GOAL := help

OAPI_CODEGEN_VERSION ?= $(shell go list -m -f '{{.Version}}' github.com/oapi-codegen/oapi-codegen/v2)
OAPI_CODEGEN_MAIN    := cmd/oapi-codegen/main.go
OAPI_CODEGEN_URL     := https://raw.githubusercontent.com/oapi-codegen/oapi-codegen/$(OAPI_CODEGEN_VERSION)/cmd/oapi-codegen/oapi-codegen.go

.PHONY: help
help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run tests.
	go test ./...

.PHONY: build
build: ## Build the project.
	go build ./...

.PHONY: sync-oapi-codegen
sync-oapi-codegen: ## Sync the oapi-codegen CLI entrypoint with the module version.
	@echo "Syncing oapi-codegen $(OAPI_CODEGEN_VERSION)"
	@curl --fail --silent --show-error --location \
		"$(OAPI_CODEGEN_URL)" \
		--output "$(OAPI_CODEGEN_MAIN)"
