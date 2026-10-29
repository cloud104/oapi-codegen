package security

import (
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

type APIKeyConfig struct {
	Name        string
	Description string
	KeyName     string
	In          string
	Scopes      []string
}

func APIKeyForOperation(op *codegen.OperationDefinition) []*APIKeyConfig {
	return APIKeyForOperations([]codegen.OperationDefinition{*op})
}

func APIKeyForOperations(ops []codegen.OperationDefinition) []*APIKeyConfig {
	var configs []*APIKeyConfig

	for _, auth := range SchemesForOperations(ops) {
		scheme := auth.Scheme
		if scheme.Type != "apiKey" {
			continue
		}

		configs = append(configs, &APIKeyConfig{
			Name:        auth.Definition.ProviderName,
			Description: scheme.Description,
			KeyName:     scheme.Name,
			In:          scheme.In,
			Scopes:      auth.Definition.Scopes,
		})
	}

	return configs
}
