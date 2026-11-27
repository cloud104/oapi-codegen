package main

import (
	"strings"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/jinzhu/inflection"
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

	codegen.TemplateFunctions["tagDescription"] = func(name string) string {
		spec, err := codegenstate.InlineSpec()
		if err != nil {
			panic(err)
		}
		for _, tag := range spec.Tags {
			if tag != nil && tag.Name == name {
				return strings.TrimSpace(tag.Description)
			}
		}
		return ""
	}

	// OpenAPI version and schema compatibility
	codegen.TemplateFunctions["isJSONSchemaCompatible"] = codegenstate.IsJSONSchemaCompatible

	// Naming
	initialisms := codegenstate.Initialisms()
	caser := naming.NewCaser(initialisms)
	codegen.TemplateFunctions["goCamelCase"] = caser.Camel
	codegen.TemplateFunctions["goPascalCase"] = caser.Pascal
	codegen.TemplateFunctions["plural"] = inflection.Plural
	codegen.TemplateFunctions["singular"] = inflection.Singular

	// Request / Responses bodies
	codegen.TemplateFunctions["jsonRequestBody"] = request.FindJSONBody
	codegen.TemplateFunctions["jsonRequestSchema"] = request.FindJSONSchema
	codegen.TemplateFunctions["genJSONRequestBodyArg"] = request.GenerateJSONBodyArgument
	codegen.TemplateFunctions["groupResponses"] = response.GroupResponses
	codegen.TemplateFunctions["bodyLimit"] = request.FindBodyLimit

	// Security / authentication
	codegen.TemplateFunctions["securityForOperation"] = security.ForOperation
	codegen.TemplateFunctions["securityForOperations"] = security.ForOperations
}
