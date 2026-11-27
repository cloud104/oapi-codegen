package codegenstate

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestSecurity(t *testing.T) {
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
	got, ok := Security()

	// Assert
	if !ok {
		t.Fatal("expected Security to return true")
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 security requirement, got %d", len(got))
	}

	scopes, ok := got[0]["BearerAuth"]
	if !ok {
		t.Fatal("expected security requirements to contain BearerAuth")
	}

	if len(scopes) != 0 {
		t.Errorf("expected no scopes, got %v", scopes)
	}
}
