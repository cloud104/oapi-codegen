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

func describeSecurityDefinition(op codegen.OperationDefinition) []codegen.SecurityDefinition {
	if op.Spec == nil {
		return nil
	}

	// Use global security unless overridden by the operation.
	security, _ := codegenstate.Security()
	if op.Spec.Security != nil {
		security = *op.Spec.Security
	}

	outDefs := make([]codegen.SecurityDefinition, 0)

	for _, sr := range security {
		for _, k := range codegen.SortedMapKeys(sr) {
			v := sr[k]
			outDefs = append(outDefs, codegen.SecurityDefinition{ProviderName: k, Scopes: v})
		}
	}

	return outDefs
}

func schemesForOperations(ops []codegen.OperationDefinition) []*OperationScheme {
	components, ok := codegenstate.Components()
	if !ok {
		return nil
	}

	scopesByProvider := make(map[string][]string)

	// Aggregate scopes by security provider.
	for _, op := range ops {
		for _, def := range describeSecurityDefinition(op) {
			provider := def.ProviderName
			scopesByProvider[provider] = append(scopesByProvider[provider], def.Scopes...)
		}
	}

	// Resolve security schemes and deduplicate scopes.
	var schemes []*OperationScheme
	for name, scopes := range scopesByProvider {
		ref := components.SecuritySchemes[name]
		if ref == nil || ref.Value == nil {
			continue
		}

		slices.Sort(scopes)
		schemes = append(schemes, &OperationScheme{
			Definition: codegen.SecurityDefinition{
				ProviderName: name,
				Scopes:       slices.Compact(scopes),
			},
			Scheme: ref.Value,
		})
	}

	// Ensure deterministic ordering.
	slices.SortFunc(schemes, func(a, b *OperationScheme) int {
		return cmp.Compare(a.Definition.ProviderName, b.Definition.ProviderName)
	})

	return schemes
}
