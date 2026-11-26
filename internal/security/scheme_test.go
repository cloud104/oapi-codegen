package security

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestSchemesForOperations(t *testing.T) {
	tests := []struct {
		name       string
		global     openapi3.SecurityRequirements
		operations func([]codegen.OperationDefinition) []codegen.OperationDefinition
		want       []codegen.SecurityDefinition
	}{
		{
			name: "all security schemes",
			want: []codegen.SecurityDefinition{
				{ProviderName: "ApiKeyAuth"},
				{ProviderName: "BasicAuth"},
				{ProviderName: "BearerAuth"},
				{ProviderName: "JWTAuth"},
			},
		},
		{
			name: "no operations",
			operations: func(_ []codegen.OperationDefinition) []codegen.OperationDefinition {
				return nil
			},
		},
		{
			name: "nil specification",
			operations: func(_ []codegen.OperationDefinition) []codegen.OperationDefinition {
				return []codegen.OperationDefinition{
					{Spec: nil},
				}
			},
		},
		{
			name: "nil security",
			operations: func(_ []codegen.OperationDefinition) []codegen.OperationDefinition {
				return []codegen.OperationDefinition{
					{Spec: &openapi3.Operation{}},
				}
			},
		},
		{
			name: "empty security",
			operations: func(_ []codegen.OperationDefinition) []codegen.OperationDefinition {
				return []codegen.OperationDefinition{
					{
						Spec: &openapi3.Operation{
							Security: &openapi3.SecurityRequirements{},
						},
					},
				}
			},
		},
		{
			name: "unknown provider",
			operations: func(_ []codegen.OperationDefinition) []codegen.OperationDefinition {
				return []codegen.OperationDefinition{
					{
						Spec: &openapi3.Operation{
							Security: &openapi3.SecurityRequirements{
								{"UnknownAuth": {}},
								{"ApiKeyAuth": {}},
							},
						},
					},
				}
			},
			want: []codegen.SecurityDefinition{
				{ProviderName: "ApiKeyAuth"},
			},
		},
		{
			name: "duplicate providers",
			operations: func(ops []codegen.OperationDefinition) []codegen.OperationDefinition {
				return append(ops, ops...)
			},
			want: []codegen.SecurityDefinition{
				{ProviderName: "ApiKeyAuth"},
				{ProviderName: "BasicAuth"},
				{ProviderName: "BearerAuth"},
				{ProviderName: "JWTAuth"},
			},
		},
		{
			name: "duplicate scopes",
			operations: func(_ []codegen.OperationDefinition) []codegen.OperationDefinition {
				return []codegen.OperationDefinition{
					{
						Spec: &openapi3.Operation{
							Security: &openapi3.SecurityRequirements{
								{"BearerAuth": {"write", "read", "read"}},
							},
						},
					},
					{
						Spec: &openapi3.Operation{
							Security: &openapi3.SecurityRequirements{
								{"BearerAuth": {"delete", "write"}},
							},
						},
					},
				}
			},
			want: []codegen.SecurityDefinition{
				{
					ProviderName: "BearerAuth",
					Scopes:       []string{"delete", "read", "write"},
				},
			},
		},
		{
			name: "multiple providers in same requirement",
			operations: func(_ []codegen.OperationDefinition) []codegen.OperationDefinition {
				return []codegen.OperationDefinition{
					{
						Spec: &openapi3.Operation{
							Security: &openapi3.SecurityRequirements{
								{
									"BearerAuth": {},
									"ApiKeyAuth": {},
									"BasicAuth":  {},
								},
							},
						},
					},
				}
			},
			want: []codegen.SecurityDefinition{
				{ProviderName: "ApiKeyAuth"},
				{ProviderName: "BasicAuth"},
				{ProviderName: "BearerAuth"},
			},
		},
		{
			name: "global security",
			global: openapi3.SecurityRequirements{
				{"BearerAuth": {}},
				{"BasicAuth": {}},
			},
			operations: func(_ []codegen.OperationDefinition) []codegen.OperationDefinition {
				return nil
			},
			want: []codegen.SecurityDefinition{
				{ProviderName: "BasicAuth"},
				{ProviderName: "BearerAuth"},
			},
		},
		{
			name: "operation security overrides global security",
			global: openapi3.SecurityRequirements{
				{"BearerAuth": {"read"}},
				{"BasicAuth": {}},
			},
			operations: func(_ []codegen.OperationDefinition) []codegen.OperationDefinition {
				return []codegen.OperationDefinition{
					{
						Spec: &openapi3.Operation{
							Security: &openapi3.SecurityRequirements{
								{"ApiKeyAuth": {}},
								{"BearerAuth": {"write"}},
							},
						},
					},
				}
			},
			want: []codegen.SecurityDefinition{
				{ProviderName: "ApiKeyAuth"},
				{
					ProviderName: "BearerAuth",
					Scopes:       []string{"write"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			doc, err := openapi3.NewLoader().LoadFromFile(
				filepath.Join("testdata", "openapi.yaml"),
			)
			if err != nil {
				t.Fatal(err)
			}

			doc.Security = tt.global

			config := codegen.Configuration{
				PackageName: "testapi",
				Generate: codegen.GenerateOptions{
					Models:    true,
					ChiServer: true,
				},
			}

			if _, err := codegen.Generate(doc, config); err != nil {
				t.Fatal(err)
			}

			operations, err := codegen.OperationDefinitions(doc)
			if err != nil {
				t.Fatal(err)
			}

			if tt.operations != nil {
				operations = tt.operations(operations)
			}

			// Act
			schemes := schemesForOperations(operations)

			// Assert
			if len(schemes) != len(tt.want) {
				t.Fatalf("expected %d security schemes, got %d", len(tt.want), len(schemes))
			}

			for i, want := range tt.want {
				got := schemes[i]

				if got == nil {
					t.Fatalf("expected non-nil scheme at index %d", i)
				}

				if got.Definition.ProviderName != want.ProviderName {
					t.Errorf("expected provider %q, got %q", want.ProviderName, got.Definition.ProviderName)
				}

				if !slices.Equal(got.Definition.Scopes, want.Scopes) {
					t.Errorf("expected scopes %v, got %v", want.Scopes, got.Definition.Scopes)
				}

				ref := doc.Components.SecuritySchemes[want.ProviderName]
				if ref == nil || ref.Value == nil {
					t.Fatalf("missing security scheme for %q", want.ProviderName)
				}

				if got.Scheme != ref.Value {
					t.Errorf("unexpected security scheme for %q", want.ProviderName)
				}
			}
		})
	}
}
