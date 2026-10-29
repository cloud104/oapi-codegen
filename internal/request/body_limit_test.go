package request

import (
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func TestFindBodyLimit(t *testing.T) {
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

	for _, op := range operations {
		if op.OperationId != "CreateItem" {
			continue
		}

		limit := FindBodyLimit(op)

		var expected int64 = 5_242_880
		if got := limit.Bytes(); got != expected {
			t.Errorf("expected body limit %d bytes, got %d", expected, got)
		}

		return
	}

	t.Fatal("operation CreateItem not found")
}
