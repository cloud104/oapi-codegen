package security

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestForOperations(t *testing.T) {
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
	s := ForOperations(operations)

	// Assert
	if !s.HasAPIKey() {
		t.Error("expected API key authentication")
	}

	if !s.HasBasic() {
		t.Error("expected basic authentication")
	}

	if !s.HasBearer() {
		t.Error("expected bearer authentication")
	}

	if !s.HasSecurity() {
		t.Error("expected security")
	}
}

func TestForOperation(t *testing.T) {
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

	if len(operations) != 3 {
		t.Fatalf("expected 3 operations, got %d", len(operations))
	}

	for _, op := range operations {
		t.Run(op.OperationId, func(t *testing.T) {
			// Act
			s := ForOperation(op)

			// Assert
			hasSecurity := len(op.SecurityDefinitions) > 0

			if s.HasSecurity() != hasSecurity {
				t.Errorf("expected HasSecurity() to be %t", hasSecurity)
			}

			if got := len(s.APIKey) + len(s.Basic) + len(s.Bearer); got != len(op.SecurityDefinitions) {
				t.Errorf("expected %d security schemes, got %d",
					len(op.SecurityDefinitions), got)
			}
		})
	}
}
