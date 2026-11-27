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

	expected := []struct {
		typ    string
		goName string
		codes  statusCodes
		fields int
	}{
		{
			typ:    "[]User",
			goName: "Users",
			codes:  statusCodes{"200"},
			fields: 1,
		},
		{
			typ:    "User",
			goName: "User",
			codes:  statusCodes{"200", "201"},
			fields: 2,
		},
		{
			typ:    "Error",
			goName: "Error",
			codes:  statusCodes{"404"},
			fields: 2,
		},
		{
			typ:    "[]openapi_types.File",
			goName: "Files",
			codes:  statusCodes{"200"},
			fields: 2,
		},
		{
			typ:    "openapi_types.File",
			goName: "File",
			codes:  statusCodes{"200"},
			fields: 1,
		},
	}

	groups := GroupResponses(defs)

	if len(groups) != len(expected) {
		t.Fatalf(
			"expected %d response groups, got %d",
			len(expected), len(groups),
		)
	}

	byType := make(map[string]*ResponseGroup, len(groups))

	for _, group := range groups {
		byType[group.Type] = group
	}

	for _, want := range expected {
		t.Run(want.goName, func(t *testing.T) {
			got, ok := byType[want.typ]
			if !ok {
				t.Fatalf("response group %q not found", want.typ)
			}

			if got.Type != want.typ {
				t.Errorf("expected type %q, got %q", want.typ, got.Type)
			}

			if got.GoName != want.goName {
				t.Errorf("expected GoName %q, got %q", want.goName, got.GoName)
			}

			if !reflect.DeepEqual(got.StatusCodes, want.codes) {
				t.Errorf("expected status codes %v, got %v", want.codes, got.StatusCodes)
			}

			if len(got.Fields) != want.fields {
				t.Errorf("expected %d fields, got %d", want.fields, len(got.Fields))
			}
		})
	}
}
