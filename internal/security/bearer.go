package security

import (
	"slices"
	"strings"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

type BearerConfig struct {
	Name         string
	Description  string
	BearerFormat string
	Scopes       []string
}

func BearerForOperation(operation *codegen.OperationDefinition) []*BearerConfig {
	return BearerForOperations([]codegen.OperationDefinition{*operation})
}

func BearerForOperations(operations []codegen.OperationDefinition) []*BearerConfig {
	var configs []*BearerConfig

	for _, auth := range SchemesForOperations(operations) {
		scheme := auth.Scheme
		if scheme.Type != "http" || !strings.EqualFold(scheme.Scheme, "bearer") {
			continue
		}

		configs = append(configs, &BearerConfig{
			Name:         auth.Definition.ProviderName,
			Description:  scheme.Description,
			BearerFormat: scheme.BearerFormat,
			Scopes:       auth.Definition.Scopes,
		})
	}

	slices.SortFunc(configs, func(a, b *BearerConfig) int {
		return strings.Compare(a.Name, b.Name)
	})

	return configs
}
