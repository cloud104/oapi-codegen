package request

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"

	"github.com/cloud104/oapi-codegen/internal/codegenstate"
)

func FindJSONBody(op *codegen.OperationDefinition) *codegen.RequestBodyDefinition {
	var body *codegen.RequestBodyDefinition

	for i := range op.Bodies {
		candidate := &op.Bodies[i]

		if !candidate.IsJSON() {
			continue
		}

		if candidate.Default {
			return candidate
		}

		if body == nil {
			body = candidate
		}
	}

	return body
}

func GenerateJSONBodyArgument(op *codegen.OperationDefinition) string {
	body := FindJSONBody(op)
	if body == nil {
		return ""
	}

	typeName := body.TypeDef(op.OperationId).TypeName

	if body.Required {
		return fmt.Sprintf(", body %s", typeName)
	}

	return fmt.Sprintf(", body *%s", typeName)
}

func FindJSONSchema(op *codegen.OperationDefinition) *openapi3.Schema {
	operations, err := codegenstate.Operations()
	if err != nil {
		panic(err)
	}

	for i := range operations {
		bundledOp := &operations[i]
		if bundledOp.Path != op.Path || !strings.EqualFold(bundledOp.Method, op.Method) {
			continue
		}

		body := FindJSONBody(bundledOp)
		if body == nil {
			continue
		}

		return body.Schema.OAPISchema
	}

	return &openapi3.Schema{}
}
