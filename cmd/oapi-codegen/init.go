package main

import (
	"fmt"
	"strings"
	_ "unsafe"
	sprig "github.com/Masterminds/sprig/v3"
	"github.com/ettle/strcase"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"strconv"
	"github.com/jinzhu/inflection"
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

	codegen.TemplateFunctions["camelCaseWithInitialisms"] = camelCaseWithInitialisms
	codegen.TemplateFunctions["genJSONRequestBodyArg"] = genJSONRequestBodyArg
	codegen.TemplateFunctions["groupResponses"] = groupResponses
	codegen.TemplateFunctions["jsonRequestBody"] = jsonRequestBody
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

func genJSONRequestBodyArg(op *codegen.OperationDefinition) string {
	body := jsonRequestBody(op)
	if body == nil {
		return ""
	}

	typeName := body.TypeDef(op.OperationId).TypeName

	return fmt.Sprintf(", body *%s", typeName)
}

type responseCodes []int

func (rc responseCodes) String() string {
	parts := make([]string, len(rc))
	for i, code := range rc {
		parts[i] = strconv.Itoa(code)
	}

	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	default:
		return strings.Join(parts[:len(parts)-1], ", ") + " or " + parts[len(parts)-1]
	}
}

type responseGroup struct {
	Type          string
	Method        string
	Fields        []string
	ResponseCodes responseCodes
}

func groupResponses(types []codegen.ResponseTypeDefinition) []responseGroup {
	groups := make(map[string]responseGroup)

	for _, def := range types {
		code, err := strconv.Atoi(def.ResponseName)
		if err != nil {
			panic(err)
		}

		typeDecl := def.Schema.TypeDecl()

		group, exists := groups[typeDecl]
		if !exists {
			method := typeDecl
			if def.Schema.ArrayType != nil {
				method = inflection.Plural(def.Schema.ArrayType.TypeDecl())
			}

			group = responseGroup{
				Type:   typeDecl,
				Method: "Get" + method,
			}
		}

		group.Fields = append(group.Fields, def.TypeName)
		group.ResponseCodes = append(group.ResponseCodes, code)

		groups[typeDecl] = group
	}

	result := make([]responseGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
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
