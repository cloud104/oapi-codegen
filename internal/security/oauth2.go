package security

import (
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

type OAuth2Config struct {
	Name        string
	Description string
	Flow        string
	Scopes      []string
	TokenURL    string
}

func OAuth2ForOperation(op *codegen.OperationDefinition) []*OAuth2Config {
	return OAuth2ForOperations([]codegen.OperationDefinition{*op})
}

func OAuth2ForOperations(ops []codegen.OperationDefinition) []*OAuth2Config {
	var configs []*OAuth2Config

	for _, security := range SchemesForOperations(ops) {
		if security.Scheme.Type != "oauth2" ||
			security.Scheme.Flows == nil {
			continue
		}

		flows := security.Scheme.Flows

		if flow := flows.ClientCredentials; flow != nil {
			configs = append(configs, &OAuth2Config{
				Name:        security.Definition.ProviderName,
				Description: security.Scheme.Description,
				Flow:        "clientCredentials",
				Scopes:      security.Definition.Scopes,
				TokenURL:    flow.TokenURL,
			})
		}

		if flow := flows.Password; flow != nil {
			configs = append(configs, &OAuth2Config{
				Name:        security.Definition.ProviderName,
				Description: security.Scheme.Description,
				Flow:        "password",
				Scopes:      security.Definition.Scopes,
				TokenURL:    flow.TokenURL,
			})
		}
	}

	return configs
}
