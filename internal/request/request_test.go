package request

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestFindJSONBody(t *testing.T) {
	doc, err := openapi3.NewLoader().LoadFromFile(
		filepath.Join("testdata", "openapi.yaml"),
	)
	if err != nil {
		t.Fatal(err)
	}

	operations, err := codegen.OperationDefinitions(doc)
	if err != nil {
		t.Fatal(err)
	}

	for i := range operations {
		op := &operations[i]

		switch op.OperationId {
		case "CreateItem":
			body := FindJSONBody(op)
			if body == nil {
				t.Fatal("expected JSON body, got nil")
			}

			if !body.IsJSON() {
				t.Errorf("expected JSON body, got %q", body.ContentType)
			}

			if !body.Required {
				t.Error("expected required JSON body")
			}

		case "UploadFile":
			if body := FindJSONBody(op); body != nil {
				t.Errorf("expected nil JSON body, got %+v", body)
			}
		}
	}
}

func TestFindJSONSchema(t *testing.T) {
	doc, err := openapi3.NewLoader().LoadFromFile(
		filepath.Join("testdata", "openapi.yaml"),
	)
	if err != nil {
		t.Fatal(err)
	}

	config := codegen.Configuration{
		PackageName: "testapi",
		Generate: codegen.GenerateOptions{
			ChiServer: true,
			Models:    true,
		},
	}

	if _, err := codegen.Generate(doc, config); err != nil {
		t.Fatal(err)
	}

	operations, err := codegen.OperationDefinitions(doc)
	if err != nil {
		t.Fatal(err)
	}

	for i := range operations {
		op := &operations[i]

		switch op.OperationId {
		case "CreateItem":
			schema := FindJSONSchema(op)
			if schema == nil {
				t.Fatal("expected JSON schema, got nil")
			}

			name, ok := schema.Properties["name"]
			if !ok {
				t.Fatal("expected name property")
			}

			if name.Value == nil || !name.Value.Type.Is("string") {
				t.Errorf("expected name property to be a string")
			}

			if len(schema.Required) != 1 || schema.Required[0] != "name" {
				t.Errorf("expected required property name, got %v", schema.Required)
			}

		case "UploadFile":
			schema := FindJSONSchema(op)
			if schema == nil {
				t.Fatal("expected empty schema, got nil")
			}

			if len(schema.Properties) != 0 {
				t.Errorf("expected empty schema, got %v", schema.Properties)
			}
		}
	}
}

func TestGenerateJSONBodyArgument(t *testing.T) {
	doc, err := openapi3.NewLoader().LoadFromFile(
		filepath.Join("testdata", "openapi.yaml"),
	)
	if err != nil {
		t.Fatal(err)
	}

	operations, err := codegen.OperationDefinitions(doc)
	if err != nil {
		t.Fatal(err)
	}

	for i := range operations {
		op := &operations[i]

		switch op.OperationId {
		case "CreateItem", "CreateOptionalItem":
			body := FindJSONBody(op)
			if body == nil {
				t.Fatal("expected JSON body, got nil")
			}

			typeName := body.TypeDef(op.OperationId).TypeName
			expected := fmt.Sprintf(", body *%s", typeName)

			if body.Required {
				expected = fmt.Sprintf(", body %s", typeName)
			}

			if got := GenerateJSONBodyArgument(op); got != expected {
				t.Errorf("expected %q, got %q", expected, got)
			}

		case "UploadFile":
			if got := GenerateJSONBodyArgument(op); got != "" {
				t.Errorf("expected empty argument, got %q", got)
			}
		}
	}
}
