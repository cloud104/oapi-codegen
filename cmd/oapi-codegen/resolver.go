package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/pb33f/libopenapi/bundler"
	"github.com/pb33f/libopenapi/datamodel"
)

func jsonRequestBodySchema(op *codegen.OperationDefinition) *openapi3.Schema {
	if op == nil {
		panic("cannot extract JSON request body schema: operation definition is nil")
	}

	specJSON, err := json.Marshal(codegenGlobalState.spec)
	if err != nil {
		panic(fmt.Errorf("failed to marshal OpenAPI spec: %w", err))
	}

	// Inline local references, including nested schema properties.
	bundled, err := bundler.BundleBytes(specJSON, &datamodel.DocumentConfiguration{
		AllowFileReferences: true,
		BundleInlineRefs:    true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to bundle OpenAPI spec: %w", err))
	}

	// Load the bundled specification into a new kin-openapi document.
	spec, err := openapi3.NewLoader().LoadFromData(bundled)
	if err != nil {
		panic(fmt.Errorf("failed to parse bundled OpenAPI spec: %w", err))
	}

	// Extract operations using oapi-codegen.
	operations, err := codegen.OperationDefinitions(spec)
	if err != nil {
		panic(fmt.Errorf("failed to extract operations from bundled OpenAPI spec: %w", err))
	}

	for i := range operations {
		bundledOp := &operations[i]
		if bundledOp.Path != op.Path || !strings.EqualFold(bundledOp.Method, op.Method) {
			continue
		}

		body := jsonRequestBody(bundledOp)
		if body == nil {
			continue
		}

		return body.Schema.OAPISchema
	}

	return &openapi3.Schema{}
}
