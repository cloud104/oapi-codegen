package main

import (
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	sprig "github.com/Masterminds/sprig/v3"
	humanize "github.com/dustin/go-humanize"
	"github.com/ettle/strcase"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jinzhu/inflection"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"

	"github.com/cloud104/oapi-codegen/internal/codegenstate"
	"github.com/cloud104/oapi-codegen/internal/security"
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

	// Naming
	codegen.TemplateFunctions["camelCaseWithInitialisms"] = camelCaseWithInitialisms
	codegen.TemplateFunctions["pascalCaseWithInitialisms"] = pascalCaseWithInitialisms

	// Request / Responses bodies
	codegen.TemplateFunctions["jsonRequestBody"] = jsonRequestBody
	codegen.TemplateFunctions["jsonRequestSchema"] = jsonRequestSchema
	codegen.TemplateFunctions["genJSONRequestBodyArg"] = genJSONRequestBodyArg
	codegen.TemplateFunctions["groupResponses"] = groupResponses
	codegen.TemplateFunctions["bodyLimit"] = bodyLimit

	// Security / authentication
	codegen.TemplateFunctions["basicAuthForOperation"] = security.BasicForOperation
	codegen.TemplateFunctions["basicAuthForOperations"] = security.BasicForOperations
	codegen.TemplateFunctions["bearerAuthForOperation"] = security.BearerForOperation
	codegen.TemplateFunctions["bearerAuthForOperations"] = security.BearerForOperations
	codegen.TemplateFunctions["apiKeyAuthForOperation"] = security.APIKeyForOperation
	codegen.TemplateFunctions["apiKeyAuthForOperations"] = security.APIKeyForOperations
}

func camelCaseWithInitialisms(s string) string {
	overrides := make(map[string]bool)

	// Surround each known initialism with separators so the caser can
	// recognize it as an independent word, even when acronyms are adjacent.
	for _, initialism := range codegenstate.Initialisms() {
		s = strings.ReplaceAll(s, initialism, "_"+initialism+"_")
		overrides[initialism] = true
	}

	return strcase.NewCaser(true, overrides, nil).ToCamel(s)
}

func pascalCaseWithInitialisms(s string) string {
	overrides := make(map[string]bool)

	// Surround each known initialism with separators so the caser can
	// recognize it as an independent word, even when acronyms are adjacent.
	for _, initialism := range codegenstate.Initialisms() {
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

func jsonRequestSchema(op *codegen.OperationDefinition) *openapi3.Schema {
	operations, err := codegenstate.Operations()
	if err != nil {
		panic(err)
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

const (
	defaultMaxBodySize = humanize.MiByte
	maxBodySize        = math.MaxInt64

	maxBodySizeExtension = "x-max-body-size"
)

type BodyLimit struct {
	bytes int64
}

func bodyLimit(op codegen.OperationDefinition) *BodyLimit {
	if op.Spec == nil || op.Spec.RequestBody == nil || op.Spec.RequestBody.Value == nil {
		return defaultBodyLimit()
	}

	raw, ok := op.Spec.RequestBody.Value.Extensions[maxBodySizeExtension]
	if !ok {
		return defaultBodyLimit()
	}

	var value string

	switch v := raw.(type) {
	case float64:
		value = strconv.FormatFloat(v, 'f', -1, 64)
	default:
		value = fmt.Sprint(v)
	}

	limit, err := parseBodyLimit(value)
	if err != nil {
		panic(fmt.Errorf("invalid body limit: %w", err))
	}

	return limit
}

func (l *BodyLimit) Bytes() int64 {
	return l.bytes
}

func (l *BodyLimit) String() string {
	return humanize.IBytes(uint64(l.bytes))
}

func defaultBodyLimit() *BodyLimit {
	return &BodyLimit{bytes: defaultMaxBodySize}
}

func newBodyLimit(n int64) (*BodyLimit, error) {
	if n < 1 || n > maxBodySize {
		return nil, fmt.Errorf("body limit must be between 1 and %d bytes", maxBodySize)
	}

	return &BodyLimit{bytes: n}, nil
}

func newBodyLimitFromUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](n T) (*BodyLimit, error) {
	if uint64(n) > math.MaxInt64 {
		return nil, fmt.Errorf("body limit exceeds int64 range")
	}

	return newBodyLimit(int64(n))
}

func parseBodyLimit(s string) (*BodyLimit, error) {
	n, err := humanize.ParseBytes(strings.TrimSpace(s))
	if err != nil {
		return nil, fmt.Errorf("parse body limit: %w", err)
	}

	return newBodyLimitFromUint(n)
}
