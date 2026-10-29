package response

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestGroupResponses(t *testing.T) {
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

	var defs []codegen.ResponseTypeDefinition

	for _, operation := range operations {
		responseTypes, err := operation.GetResponseTypeDefinitions()
		if err != nil {
			t.Fatal(err)
		}

		defs = append(defs, responseTypes...)
	}

	groups := GroupResponses(defs)
	if len(groups) != 3 {
		t.Fatalf("expected 3 response groups, got %d", len(groups))
	}

	expected := []struct {
		typ    string
		method string
		codes  statusCodes
		fields int
	}{
		{
			typ:    "[]User",
			method: "GetUsers",
			codes:  statusCodes{"200"},
			fields: 1,
		},
		{
			typ:    "User",
			method: "GetUser",
			codes:  statusCodes{"200", "201"},
			fields: 2,
		},
		{
			typ:    "Error",
			method: "GetError",
			codes:  statusCodes{"404"},
			fields: 2,
		},
	}
	for i, want := range expected {
		got := groups[i]
		if got.Type != want.typ {
			t.Errorf("group %d: expected type %q, got %q", i, want.typ, got.Type)
		}

		if got.Method != want.method {
			t.Errorf("group %d: expected method %q, got %q", i, want.method, got.Method)
		}

		if !reflect.DeepEqual(got.StatusCodes, want.codes) {
			t.Errorf("group %d: expected status codes %v, got %v", i, want.codes, got.StatusCodes)
		}

		if len(got.Fields) != want.fields {
			t.Errorf("group %d: expected %d fields, got %d", i, want.fields, len(got.Fields))
		}
	}
}
