package codegenstate

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestComponents(t *testing.T) {
	// Arrange
	loader := openapi3.NewLoader()

	spec, err := loader.LoadFromFile(
		filepath.Join("testdata", "openapi.yaml"),
	)
	if err != nil {
		t.Fatalf("failed to load OpenAPI spec: %v", err)
	}

	config := codegen.Configuration{
		PackageName: "testapi",
		Generate: codegen.GenerateOptions{
			ChiServer: true,
			Models:    true,
		},
	}

	_, err = codegen.Generate(spec, config)
	if err != nil {
		t.Fatalf("failed to generate code: %v", err)
	}

	// Act
	got, ok := Components()

	// Assert
	if !ok {
		t.Fatal("expected Components to return true")
	}

	if got == nil {
		t.Fatal("expected components, got nil")
	}

	if _, ok := got.Schemas["User"]; !ok {
		t.Error("expected components to contain User schema")
	}
}
