package naming

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestSingular(t *testing.T) {
	configureInitialisms(t)

	tests := []struct {
		input string
		want  string
	}{
		// Regular plurals.
		{"users", "user"},
		{"books", "book"},
		{"categories", "category"},
		{"boxes", "box"},
		{"statuses", "status"},

		// Irregular plurals.
		{"people", "person"},
		{"children", "child"},
		{"men", "man"},
		{"mice", "mouse"},

		// Already singular.
		{"user", "user"},
		{"person", "person"},

		// Uncountable nouns.
		{"sheep", "sheep"},
		{"fish", "fish"},
		{"equipment", "equipment"},

		// Uppercase initialisms.
		{"TKS", "TKS"},
		{"AWS", "AWS"},
		{"DNS", "DNS"},
		{"API", "API"},
		{"URL", "URL"},

		// Lowercase initialisms.
		{"tks", "tks"},
		{"aws", "aws"},
		{"dns", "dns"},
		{"api", "api"},
		{"url", "url"},

		// Mixed-case initialisms.
		{"Tks", "Tks"},
		{"Aws", "Aws"},
		{"Dns", "Dns"},
		{"Api", "Api"},
		{"Url", "Url"},

		// Uppercase initialisms with a plural suffix.
		{"TKSs", "TKS"},
		{"AWSs", "AWS"},
		{"DNSs", "DNS"},
		{"APIs", "API"},
		{"URLs", "URL"},

		// Lowercase plural initialisms.
		{"tkss", "tks"},
		{"awss", "aws"},
		{"dnss", "dns"},
		{"apis", "api"},
		{"urls", "url"},

		// Mixed-case plural initialisms.
		{"Tkss", "Tks"},
		{"Awss", "Aws"},
		{"Dnss", "Dns"},
		{"Apis", "Api"},
		{"Urls", "Url"},

		// PascalCase.
		{"Users", "User"},
		{"Categories", "Category"},

		// Uppercase nouns that are not registered initialisms.
		{"USER", "USER"},
		{"USERS", "USER"},

		// Empty input.
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Singular(tt.input); got != tt.want {
				t.Errorf("Singular(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPlural(t *testing.T) {
	configureInitialisms(t)

	tests := []struct {
		input string
		want  string
	}{
		// Regular nouns.
		{"user", "users"},
		{"book", "books"},
		{"category", "categories"},
		{"box", "boxes"},
		{"status", "statuses"},

		// Irregular nouns.
		{"person", "people"},
		{"child", "children"},
		{"man", "men"},
		{"mouse", "mice"},

		// Already plural.
		{"users", "users"},
		{"people", "people"},
		{"children", "children"},
		{"men", "men"},
		{"mice", "mice"},

		// Uncountable nouns.
		{"sheep", "sheep"},
		{"fish", "fish"},
		{"equipment", "equipment"},

		// Uppercase initialisms.
		{"TKS", "TKS"},
		{"AWS", "AWS"},
		{"DNS", "DNS"},
		{"API", "API"},
		{"URL", "URL"},

		// Lowercase initialisms.
		{"tks", "tks"},
		{"aws", "aws"},
		{"dns", "dns"},
		{"api", "api"},
		{"url", "url"},

		// Mixed-case initialisms.
		{"Tks", "Tks"},
		{"Aws", "Aws"},
		{"Dns", "Dns"},
		{"Api", "Api"},
		{"Url", "Url"},

		// Uppercase initialisms with an existing plural suffix.
		{"TKSs", "TKSs"},
		{"AWSs", "AWSs"},
		{"DNSs", "DNSs"},
		{"APIs", "APIs"},
		{"URLs", "URLs"},

		// Already plural lowercase initialisms.
		{"tkss", "tkss"},
		{"awss", "awss"},
		{"dnss", "dnss"},
		{"apis", "apis"},
		{"urls", "urls"},

		// Already plural mixed-case initialisms.
		{"Tkss", "Tkss"},
		{"Awss", "Awss"},
		{"Dnss", "Dnss"},
		{"Apis", "Apis"},
		{"Urls", "Urls"},

		// PascalCase.
		{"User", "Users"},
		{"Category", "Categories"},

		// Uppercase nouns that are not registered initialisms.
		{"USER", "USERS"},
		{"USERS", "USERS"},

		// Empty input.
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Plural(tt.input); got != tt.want {
				t.Errorf("Plural(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func configureInitialisms(t *testing.T) {
	t.Helper()

	spec, err := openapi3.NewLoader().LoadFromData([]byte(`
openapi: 3.0.3
info:
  title: Test API
  version: 1.0.0
paths: {}
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
`))
	if err != nil {
		t.Fatalf("failed to load OpenAPI specification: %v", err)
	}

	if err := spec.Validate(context.Background()); err != nil {
		t.Fatalf("failed to validate OpenAPI specification: %v", err)
	}

	config := codegen.Configuration{
		PackageName: "testapi",
		Generate: codegen.GenerateOptions{
			ChiServer: true,
			Models:    true,
		},
		OutputOptions: codegen.OutputOptions{
			NameNormalizer: "ToCamelCaseWithInitialisms",
			AdditionalInitialisms: []string{
				"TKS",
				"AWS",
				"DNS",
				"API",
				"URL",
			},
		},
	}

	_, err = codegen.Generate(spec, config)
	if err != nil {
		t.Fatalf("failed to generate code: %v", err)
	}

	// Ensure that additional initialisms were successfully registered.
	for _, initialism := range config.OutputOptions.AdditionalInitialisms {
		if !isInitialism(initialism) {
			t.Fatalf("initialism %q was not registered", initialism)
		}
	}
}
