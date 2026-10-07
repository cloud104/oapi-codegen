package main

import (
	"fmt"
	"strings"
	_ "unsafe"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/ettle/strcase"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

// Mirrors the prefix of codegen.globalState up to initialismsMap.
// This is intentionally coupled to oapi-codegen's internal implementation.
type codegenGlobalStateLayout struct {
	options        codegen.Configuration
	spec           *openapi3.T
	is31           bool
	importMapping  map[string]struct{ Name, Path string }
	initialismsMap map[string]string
}

//go:linkname codegenGlobalState github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen.globalState
var codegenGlobalState codegenGlobalStateLayout

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

	// Naming
	codegen.TemplateFunctions["camelCaseWithInitialisms"] = camelCaseWithInitialisms
	codegen.TemplateFunctions["pascalCaseWithInitialisms"] = pascalCaseWithInitialisms

	// Request / Responses bodies
	codegen.TemplateFunctions["jsonRequestBody"] = jsonRequestBody
	codegen.TemplateFunctions["genJSONRequestBodyArg"] = genJSONRequestBodyArg
}

func camelCaseWithInitialisms(s string) string {
	overrides := make(map[string]bool)

	// Surround each known initialism with separators so the caser can
	// recognize it as an independent word, even when acronyms are adjacent.
	for _, initialism := range codegenGlobalState.initialismsMap {
		s = strings.ReplaceAll(s, initialism, "_"+initialism+"_")
		overrides[initialism] = true
	}

	return strcase.NewCaser(true, overrides, nil).ToCamel(s)
}

func pascalCaseWithInitialisms(s string) string {
	overrides := make(map[string]bool)

	// Surround each known initialism with separators so the caser can
	// recognize it as an independent word, even when acronyms are adjacent.
	for _, initialism := range codegenGlobalState.initialismsMap {
		s = strings.ReplaceAll(s, initialism, "_"+initialism+"_")
		overrides[initialism] = true
	}

	return strcase.NewCaser(true, overrides, nil).ToPascal(s)
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
