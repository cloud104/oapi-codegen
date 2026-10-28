package security

import (
	"cmp"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"

	"github.com/cloud104/oapi-codegen/internal/codegenstate"
)

type OperationScheme struct {
	Definition codegen.SecurityDefinition
	Scheme     *openapi3.SecurityScheme
}

func SchemesForOperations(ops []codegen.OperationDefinition) []*OperationScheme {
	components, ok := codegenstate.Components()
	if !ok {
		return nil
	}

	seen := make(map[string]bool)

	var schemes []*OperationScheme

	for _, op := range ops {
		for _, def := range op.SecurityDefinitions {
			name := def.ProviderName

			ref := components.SecuritySchemes[name]
			if seen[name] || ref == nil || ref.Value == nil {
				continue
			}

			seen[name] = true

			schemes = append(schemes, &OperationScheme{def, ref.Value})
		}
	}

	slices.SortFunc(schemes, func(a, b *OperationScheme) int {
		return cmp.Compare(a.Definition.ProviderName, b.Definition.ProviderName)
	})

	return schemes
}
