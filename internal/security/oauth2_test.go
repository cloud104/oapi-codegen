package security

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestOAuth2ForOperations(t *testing.T) {
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
	configs := ForOperations(operations).OAuth2

	// Assert
	if len(configs) != 1 {
		t.Fatalf("expected 1 OAuth2 config, got %d", len(configs))
	}

	got := configs[0]

	if got.Name != "OAuth2Auth" {
		t.Errorf("expected name OAuth2Auth, got %q", got.Name)
	}

	if got.Description != "OAuth2 authentication" {
		t.Errorf("unexpected description %q", got.Description)
	}

	if got.Flow != "clientCredentials" {
		t.Errorf("expected flow clientCredentials, got %q", got.Flow)
	}

	if got.TokenURL != "https://example.com/oauth/token" {
		t.Errorf("unexpected token URL %q", got.TokenURL)
	}

	if want := []string{"read"}; !slices.Equal(got.Scopes, want) {
		t.Errorf("expected scopes %v, got %v", want, got.Scopes)
	}
}
