package main

import (
	"fmt"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func init() {
	for name, fn := range sprig.FuncMap() {
		if _, exists := codegen.TemplateFunctions[name]; !exists {
			codegen.TemplateFunctions[name] = fn
		}
	}

	codegen.TemplateFunctions["jsonRequestBody"] = jsonRequestBody
	codegen.TemplateFunctions["genJSONRequestBodyArg"] = genJSONRequestBodyArg
}

func jsonRequestBody(op *codegen.OperationDefinition) *codegen.RequestBodyDefinition {
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

func genJSONRequestBodyArg(op *codegen.OperationDefinition) string {
	body := jsonRequestBody(op)
	if body == nil {
		return ""
	}

	typeName := body.TypeDef(op.OperationId).TypeName

	return fmt.Sprintf(", body *%s", typeName)
}
