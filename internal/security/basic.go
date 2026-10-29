package security

import (
	"strings"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

type BasicConfig struct {
	Name        string
	Description string
	Scopes      []string
}

func BasicForOperation(op *codegen.OperationDefinition) []*BasicConfig {
	return BasicForOperations([]codegen.OperationDefinition{*op})
}

func BasicForOperations(ops []codegen.OperationDefinition) []*BasicConfig {
	var configs []*BasicConfig

	for _, auth := range SchemesForOperations(ops) {
		scheme := auth.Scheme
		if scheme.Type != "http" || !strings.EqualFold(scheme.Scheme, "basic") {
			continue
		}

		configs = append(configs, &BasicConfig{
			Name:        auth.Definition.ProviderName,
			Description: scheme.Description,
			Scopes:      auth.Definition.Scopes,
		})
	}

	return configs
}
