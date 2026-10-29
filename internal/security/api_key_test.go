package security

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestAPIKeyForOperations(t *testing.T) {
	// Arrange
	doc, err := openapi3.NewLoader().LoadFromFile(
		filepath.Join("testdata", "openapi.yaml"),
	)
	if err != nil {
		t.Fatal(err)
	}

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

	// Act
	configs := APIKeyForOperations(operations)

	// Assert
	if len(configs) != 1 {
		t.Fatalf("expected 1 API key config, got %d", len(configs))
	}

	got := configs[0]
	if got.Name != "ApiKeyAuth" {
		t.Errorf("expected name ApiKeyAuth, got %q", got.Name)
	}

	if got.Description != "API key authentication" {
		t.Errorf("unexpected description %q", got.Description)
	}

	if got.KeyName != "X-API-Key" {
		t.Errorf("expected key name X-API-Key, got %q", got.KeyName)
	}

	if got.In != "header" {
		t.Errorf("expected location header, got %q", got.In)
	}

	if len(got.Scopes) != 0 {
		t.Errorf("expected no scopes, got %v", got.Scopes)
	}
}
