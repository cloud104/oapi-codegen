package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	_ "unsafe"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/ettle/strcase"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jinzhu/inflection"
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
	codegen.TemplateFunctions["groupResponses"] = groupResponses
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

type responseCodes []string

func (rc responseCodes) String() string {
	parts := make([]string, len(rc))
	for i, code := range rc {
		parts[i] = code
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

var (
	structRE = regexp.MustCompile(`(?s)^\s*struct\s*\{.*\}\s*$`)
	mapRE    = regexp.MustCompile(`^\s*map\s*\[.*\].*$`)
)

func groupResponses(types []codegen.ResponseTypeDefinition) []*responseGroup {
	groups := make(map[string]*responseGroup)

	for _, def := range types {
		typeDecl := def.Schema.TypeDecl()
		if structRE.MatchString(typeDecl) || mapRE.MatchString(typeDecl) {
			continue
		}

		group, exists := groups[typeDecl]
		if !exists {
			method := typeDecl
			if def.Schema.ArrayType != nil {
				method = inflection.Plural(def.Schema.ArrayType.TypeDecl())
			}

			group = &responseGroup{
				Type:   typeDecl,
				Method: "Get" + method,
			}
			groups[typeDecl] = group
		}

		group.Fields = append(group.Fields, def.TypeName)
		if !slices.Contains(group.ResponseCodes, def.ResponseName) {
			group.ResponseCodes = append(group.ResponseCodes, def.ResponseName)
		}
	}

	result := make([]*responseGroup, 0, len(groups))
	for _, group := range groups {
		slices.Sort(group.Fields)
		slices.Sort(group.ResponseCodes)
		result = append(result, group)
	}

	// Map iteration order is undefined, so sort the groups to keep generated
	// output deterministic. Response codes define the primary order, with the
	// response type acting as a stable tie-breaker.
	slices.SortFunc(result, func(a, b *responseGroup) int {
		if cmp := slices.Compare(a.ResponseCodes, b.ResponseCodes); cmp != 0 {
			return cmp
		}

		return strings.Compare(a.Type, b.Type)
	})

	return result
}
