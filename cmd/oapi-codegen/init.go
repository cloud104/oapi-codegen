package main

import (
	"fmt"
	"strings"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/ettle/strcase"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func init() {
	for name, fn := range sprig.FuncMap() {
		exists := false

		for existingName := range codegen.TemplateFunctions {
			if strings.EqualFold(existingName, name) {
				exists = true
				break
			}
		}

		if !exists {
			codegen.TemplateFunctions[name] = fn
		}
	}

	codegen.TemplateFunctions["camelCaseWithInitialisms"] = camelCaseWithInitialisms
	codegen.TemplateFunctions["genJSONRequestBodyArg"] = genJSONRequestBodyArg
	codegen.TemplateFunctions["jsonRequestBody"] = jsonRequestBody
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

func camelCaseWithInitialisms(s string, initialisms ...string) string {
	// Surround each known initialism with separators so the caser can
	// recognize it as an independent word, even when acronyms are adjacent.
	for _, initialism := range initialisms {
		s = strings.ReplaceAll(s, initialism, "_"+initialism+"_")
	}

	initialismOverrides := make(map[string]bool, len(initialisms))
	for _, initialism := range initialisms {
		initialismOverrides[initialism] = true
	}

	caser := strcase.NewCaser(true, initialismOverrides, nil)

	return caser.ToCamel(s)
}
