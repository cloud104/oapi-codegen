package security

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestBearerForOperations(t *testing.T) {
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
	configs := BearerForOperations(operations)

	// Assert
	if len(configs) != 2 {
		t.Fatalf("expected 2 bearer configs, got %d", len(configs))
	}

	expected := []struct {
		name        string
		description string
		format      string
	}{
		{
			name:        "BearerAuth",
			description: "Bearer authentication",
			format:      "opaque",
		},
		{
			name:        "JWTAuth",
			description: "JWT authentication",
			format:      "JWT",
		},
	}

	for i, want := range expected {
		got := configs[i]

		if got.Name != want.name {
			t.Errorf("config %d: expected name %q, got %q", i, want.name, got.Name)
		}

		if got.Description != want.description {
			t.Errorf("config %d: expected description %q, got %q", i, want.description, got.Description)
		}

		if got.BearerFormat != want.format {
			t.Errorf("config %d: expected format %q, got %q", i, want.format, got.BearerFormat)
		}

		if len(got.Scopes) != 0 {
			t.Errorf("config %d: expected no scopes, got %v", i, got.Scopes)
		}
	}
}
