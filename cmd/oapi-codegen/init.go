package main

import (
	"strings"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"

	"github.com/cloud104/oapi-codegen/internal/codegenstate"
	"github.com/cloud104/oapi-codegen/internal/naming"
	"github.com/cloud104/oapi-codegen/internal/request"
	"github.com/cloud104/oapi-codegen/internal/response"
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

	// Security / authentication
	codegen.TemplateFunctions["basicAuth"] = security.BasicForOperation
	codegen.TemplateFunctions["basicAuthFromList"] = security.BasicForOperations
	codegen.TemplateFunctions["bearerAuth"] = security.BearerForOperation
	codegen.TemplateFunctions["bearerAuthFromList"] = security.BearerForOperations
	codegen.TemplateFunctions["apiKeyAuth"] = security.APIKeyForOperation
	codegen.TemplateFunctions["apiKeyAuthFromList"] = security.APIKeyForOperations
}
