package codegenstate

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestInlineSpec(t *testing.T) {
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
			Models: true,
		},
	}

	// Act
	_, err = codegen.Generate(spec, config)
	if err != nil {
		t.Fatalf("failed to generate code: %v", err)
	}

	got, err := InlineSpec()
	if err != nil {
		t.Fatalf("failed to inline spec: %v", err)
	}

	// Assert
	if got == nil {
		t.Fatal("expected inlined spec, got nil")
	}

	if got.OpenAPI == "" {
		t.Error("expected inlined spec to have an OpenAPI version")
	}

	if gotPaths, specPaths := len(got.Paths.Map()), len(spec.Paths.Map()); gotPaths != specPaths {
		t.Errorf("unexpected number of paths: got %d, want %d", gotPaths, specPaths)
	}

	userRef, ok := got.Components.Schemas["User"]
	if !ok {
		t.Fatal("expected inlined spec to contain User schema")
	}

	if userRef.Value == nil {
		t.Error("expected User schema to be resolved inline")
	}
}
