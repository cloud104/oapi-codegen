package security

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestBasicForOperations(t *testing.T) {
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
	configs := BasicForOperations(operations)

	// Assert
	if len(configs) != 1 {
		t.Fatalf("expected 1 basic auth config, got %d", len(configs))
	}

	got := configs[0]

	if got.Name != "BasicAuth" {
		t.Errorf("expected name BasicAuth, got %q", got.Name)
	}

	if got.Description != "Basic authentication" {
		t.Errorf("unexpected description %q", got.Description)
	}

	if len(got.Scopes) != 0 {
		t.Errorf("expected no scopes, got %v", got.Scopes)
	}
}
