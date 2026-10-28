package codegenstate

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestOperations(t *testing.T) {
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
	got, err := Operations()
	if err != nil {
		t.Fatalf("failed to get operations: %v", err)
	}

	// Assert
	if len(got) != 2 {
		t.Fatalf("unexpected number of operations: got %d, want %d", len(got), 2)
	}

	operationsByID := make(map[string]codegen.OperationDefinition, len(got))
	for _, operation := range got {
		operationsByID[operation.OperationId] = operation
	}

	for _, operationID := range []string{"ListUsers", "GetUser"} {
		if _, ok := operationsByID[operationID]; !ok {
			t.Errorf("expected operation %q to be present", operationID)
		}
	}
}
