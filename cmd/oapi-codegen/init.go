package main

import (
	"cmp"
	"slices"
	"strings"
	"github.com/samber/lo"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"

	"github.com/cloud104/oapi-codegen/internal/codegenstate"
	"github.com/cloud104/oapi-codegen/internal/naming"
	"github.com/cloud104/oapi-codegen/internal/request"
	"github.com/cloud104/oapi-codegen/internal/response"
	"github.com/cloud104/oapi-codegen/internal/security"
	"regexp"
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
	initialisms := codegenstate.Initialisms()
	caser := naming.NewCaser(initialisms)
	codegen.TemplateFunctions["goCamelCase"] = caser.Camel
	codegen.TemplateFunctions["goPascalCase"] = caser.Pascal

	// Request / Responses bodies
	codegen.TemplateFunctions["jsonRequestBody"] = request.FindJSONBody
	codegen.TemplateFunctions["jsonRequestSchema"] = request.FindJSONSchema
	codegen.TemplateFunctions["genJSONRequestBodyArg"] = request.GenerateJSONBodyArgument
	codegen.TemplateFunctions["groupResponses"] = response.GroupResponses
	codegen.TemplateFunctions["bodyLimit"] = request.FindBodyLimit
	codegen.TemplateFunctions["getResponses"] = getResponses

	// Security / authentication
	codegen.TemplateFunctions["basicAuth"] = security.BasicForOperation
	codegen.TemplateFunctions["basicAuthFromList"] = security.BasicForOperations
	codegen.TemplateFunctions["bearerAuth"] = security.BearerForOperation
	codegen.TemplateFunctions["bearerAuthFromList"] = security.BearerForOperations
	codegen.TemplateFunctions["apiKeyAuth"] = security.APIKeyForOperation
	codegen.TemplateFunctions["apiKeyAuthFromList"] = security.APIKeyForOperations
}

var (
	structRE = regexp.MustCompile(`(?s)^\s*struct\s*\{.*\}\s*$`)
	mapRE    = regexp.MustCompile(`^\s*map\s*\[.*\].*$`)
)

type Response struct {
	Code   int
	Name   string
	GoName string
}

func getResponses(op codegen.OperationDefinition) []*Response {
	definitions, err := op.GetResponseTypeDefinitions()
	if err != nil {
		panic(err)
	}

	var definitionMap = lo.GroupBy(definitions, func(def codegen.ResponseTypeDefinition) string {
		return def.ResponseName
	})

	responseMap := make(map[string]*Response)
	for name, typeDefs := range definitionMap {
		status, ok := httpStatuses[name]
		if !ok {
			continue
		}

		responses := response.GroupResponses(typeDefs)
		_ = responses

		responseMap[name] = status
	}

	responses := lo.Values(responseMap)

	slices.SortFunc(responses, func(a, b *Response) int {
		return cmp.Compare(a.Code, b.Code)
	})

	return responses
}

// HTTP status codes as registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
var httpStatuses = map[string]*Response{
	// 1xx: Informational
	"100": {Code: 100, Name: "Continue", GoName: "http.StatusContinue"},
	"101": {Code: 101, Name: "SwitchingProtocols", GoName: "http.StatusSwitchingProtocols"},
	"102": {Code: 102, Name: "Processing", GoName: "http.StatusProcessing"},
	"103": {Code: 103, Name: "EarlyHints", GoName: "http.StatusEarlyHints"},

	// 2xx: Success
	"200": {Code: 200, Name: "OK", GoName: "http.StatusOK"},
	"201": {Code: 201, Name: "Created", GoName: "http.StatusCreated"},
	"202": {Code: 202, Name: "Accepted", GoName: "http.StatusAccepted"},
	"203": {Code: 203, Name: "NonAuthoritativeInfo", GoName: "http.StatusNonAuthoritativeInfo"},
	"204": {Code: 204, Name: "NoContent", GoName: "http.StatusNoContent"},
	"205": {Code: 205, Name: "ResetContent", GoName: "http.StatusResetContent"},
	"206": {Code: 206, Name: "PartialContent", GoName: "http.StatusPartialContent"},
	"207": {Code: 207, Name: "MultiStatus", GoName: "http.StatusMultiStatus"},
	"208": {Code: 208, Name: "AlreadyReported", GoName: "http.StatusAlreadyReported"},
	"226": {Code: 226, Name: "IMUsed", GoName: "http.StatusIMUsed"},

	// 3xx: Redirection
	"300": {Code: 300, Name: "MultipleChoices", GoName: "http.StatusMultipleChoices"},
	"301": {Code: 301, Name: "MovedPermanently", GoName: "http.StatusMovedPermanently"},
	"302": {Code: 302, Name: "Found", GoName: "http.StatusFound"},
	"303": {Code: 303, Name: "SeeOther", GoName: "http.StatusSeeOther"},
	"304": {Code: 304, Name: "NotModified", GoName: "http.StatusNotModified"},
	"305": {Code: 305, Name: "UseProxy", GoName: "http.StatusUseProxy"},
	"307": {Code: 307, Name: "TemporaryRedirect", GoName: "http.StatusTemporaryRedirect"},
	"308": {Code: 308, Name: "PermanentRedirect", GoName: "http.StatusPermanentRedirect"},

	// 4xx: Client errors
	"400": {Code: 400, Name: "BadRequest", GoName: "http.StatusBadRequest"},
	"401": {Code: 401, Name: "Unauthorized", GoName: "http.StatusUnauthorized"},
	"402": {Code: 402, Name: "PaymentRequired", GoName: "http.StatusPaymentRequired"},
	"403": {Code: 403, Name: "Forbidden", GoName: "http.StatusForbidden"},
	"404": {Code: 404, Name: "NotFound", GoName: "http.StatusNotFound"},
	"405": {Code: 405, Name: "MethodNotAllowed", GoName: "http.StatusMethodNotAllowed"},
	"406": {Code: 406, Name: "NotAcceptable", GoName: "http.StatusNotAcceptable"},
	"407": {Code: 407, Name: "ProxyAuthRequired", GoName: "http.StatusProxyAuthRequired"},
	"408": {Code: 408, Name: "RequestTimeout", GoName: "http.StatusRequestTimeout"},
	"409": {Code: 409, Name: "Conflict", GoName: "http.StatusConflict"},
	"410": {Code: 410, Name: "Gone", GoName: "http.StatusGone"},
	"411": {Code: 411, Name: "LengthRequired", GoName: "http.StatusLengthRequired"},
	"412": {Code: 412, Name: "PreconditionFailed", GoName: "http.StatusPreconditionFailed"},
	"413": {Code: 413, Name: "RequestEntityTooLarge", GoName: "http.StatusRequestEntityTooLarge"},
	"414": {Code: 414, Name: "RequestURITooLong", GoName: "http.StatusRequestURITooLong"},
	"415": {Code: 415, Name: "UnsupportedMediaType", GoName: "http.StatusUnsupportedMediaType"},
	"416": {Code: 416, Name: "RequestedRangeNotSatisfiable", GoName: "http.StatusRequestedRangeNotSatisfiable"},
	"417": {Code: 417, Name: "ExpectationFailed", GoName: "http.StatusExpectationFailed"},
	"418": {Code: 418, Name: "Teapot", GoName: "http.StatusTeapot"},
	"421": {Code: 421, Name: "MisdirectedRequest", GoName: "http.StatusMisdirectedRequest"},
	"422": {Code: 422, Name: "UnprocessableEntity", GoName: "http.StatusUnprocessableEntity"},
	"423": {Code: 423, Name: "Locked", GoName: "http.StatusLocked"},
	"424": {Code: 424, Name: "FailedDependency", GoName: "http.StatusFailedDependency"},
	"425": {Code: 425, Name: "TooEarly", GoName: "http.StatusTooEarly"},
	"426": {Code: 426, Name: "UpgradeRequired", GoName: "http.StatusUpgradeRequired"},
	"428": {Code: 428, Name: "PreconditionRequired", GoName: "http.StatusPreconditionRequired"},
	"429": {Code: 429, Name: "TooManyRequests", GoName: "http.StatusTooManyRequests"},
	"431": {Code: 431, Name: "RequestHeaderFieldsTooLarge", GoName: "http.StatusRequestHeaderFieldsTooLarge"},
	"451": {Code: 451, Name: "UnavailableForLegalReasons", GoName: "http.StatusUnavailableForLegalReasons"},

	// 5xx: Server errors
	"500": {Code: 500, Name: "InternalServerError", GoName: "http.StatusInternalServerError"},
	"501": {Code: 501, Name: "NotImplemented", GoName: "http.StatusNotImplemented"},
	"502": {Code: 502, Name: "BadGateway", GoName: "http.StatusBadGateway"},
	"503": {Code: 503, Name: "ServiceUnavailable", GoName: "http.StatusServiceUnavailable"},
	"504": {Code: 504, Name: "GatewayTimeout", GoName: "http.StatusGatewayTimeout"},
	"505": {Code: 505, Name: "HTTPVersionNotSupported", GoName: "http.StatusHTTPVersionNotSupported"},
	"506": {Code: 506, Name: "VariantAlsoNegotiates", GoName: "http.StatusVariantAlsoNegotiates"},
	"507": {Code: 507, Name: "InsufficientStorage", GoName: "http.StatusInsufficientStorage"},
	"508": {Code: 508, Name: "LoopDetected", GoName: "http.StatusLoopDetected"},
	"510": {Code: 510, Name: "NotExtended", GoName: "http.StatusNotExtended"},
	"511": {Code: 511, Name: "NetworkAuthenticationRequired", GoName: "http.StatusNetworkAuthenticationRequired"},
}

var contentTypes = map[string]string{
	// JSON
	"application/json":         "JSON",
	"text/x-json":              "JSON",
	"application/problem+json": "ProblemJSON",

	// HAL JSON
	"application/hal+json": "HALJSON",

	// YAML
	"application/yaml":   "YAML",
	"application/x-yaml": "YAML",
	"text/yaml":          "YAML",
	"text/x-yaml":        "YAML",

	// XML
	"application/xml":         "XML",
	"text/xml":                "XML",
	"application/problem+xml": "ProblemXML",
}
