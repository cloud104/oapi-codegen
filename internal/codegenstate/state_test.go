package codegenstate

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestGlobalState(t *testing.T) {
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

	// Assert
	if got := globalState.options.PackageName; got != config.PackageName {
		t.Errorf("unexpected package name: got %q, want %q", got, config.PackageName)
	}

	if globalState.spec != spec {
		t.Error("globalState.spec does not match the loaded OpenAPI spec")
	}

	if got, want := globalState.is31, spec.IsOpenAPI31OrLater(); got != want {
		t.Errorf("unexpected is31: got %t, want %t", got, want)
	}

	if got := globalState.initialismsMap["api"]; got != "API" {
		t.Errorf("unexpected initialism: got %q, want %q", got, "API")
	}
}
