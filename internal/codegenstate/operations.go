package codegenstate

import (
	"errors"
	"fmt"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

var (
	operationsOnce sync.Once
	operations     []codegen.OperationDefinition
	operationsErr  error
)

func Operations() ([]codegen.OperationDefinition, error) {
	operationsOnce.Do(func() {
		spec, err := InlineSpec()
		if err != nil {
			operationsErr = err
			return
		}

		if spec == nil {
			operationsErr = errors.New("OpenAPI spec is nil")
			return
		}

		if spec.Paths == nil {
			operationsErr = errors.New("OpenAPI spec has no paths")
			return
		}

		clearDescriptions(spec)
		clearExtensions(spec)
		clearExamples(spec)

		operations, operationsErr = codegen.OperationDefinitions(spec)
		if operationsErr != nil {
			operationsErr = fmt.Errorf("generate OpenAPI operation definitions: %w", operationsErr)
		}
	})

	return operations, operationsErr
}

func clearDescriptions(spec *openapi3.T) {
	visited := make(map[*openapi3.Schema]bool)

	for _, path := range spec.Paths.Map() {
		if path == nil {
			continue
		}

		path.Description = ""
		clearParameterDescriptions(path.Parameters, visited)

		for _, operation := range path.Operations() {
			if operation == nil {
				continue
			}

			operation.Description = ""
			clearParameterDescriptions(operation.Parameters, visited)

			if operation.Responses != nil {
				for _, responseRef := range operation.Responses.Map() {
					if responseRef == nil || responseRef.Value == nil {
						continue
					}

					response := responseRef.Value

					response.Description = nil
					for _, headerRef := range response.Headers {
						if headerRef == nil || headerRef.Value == nil {
							continue
						}

						header := headerRef.Value
						header.Description = ""
						clearSchemaDescriptions(header.Schema, visited)
						clearContentDescriptions(header.Content, visited)
					}

					clearContentDescriptions(response.Content, visited)
				}
			}

			bodyRef := operation.RequestBody
			if bodyRef == nil {
				continue
			}

			bodyRef.Description = nil
			if bodyRef.Value == nil {
				continue
			}

			body := bodyRef.Value
			body.Description = ""
			clearContentDescriptions(body.Content, visited)
		}
	}
}

func clearSchemaDescriptions(ref *openapi3.SchemaRef, visited map[*openapi3.Schema]bool) {
	if ref == nil || ref.Value == nil || visited[ref.Value] {
		return
	}

	schema := ref.Value
	visited[schema] = true

	schema.Description = ""
	for _, property := range schema.Properties {
		clearSchemaDescriptions(property, visited)
	}

	for _, refs := range [][]*openapi3.SchemaRef{
		schema.AllOf,
		schema.AnyOf,
		schema.OneOf,
	} {
		for _, ref := range refs {
			clearSchemaDescriptions(ref, visited)
		}
	}

	clearSchemaDescriptions(schema.Items, visited)
	clearSchemaDescriptions(schema.Not, visited)
	clearSchemaDescriptions(schema.AdditionalProperties.Schema, visited)
}

func clearContentDescriptions(content openapi3.Content, visited map[*openapi3.Schema]bool) {
	for _, media := range content {
		if media != nil {
			clearSchemaDescriptions(media.Schema, visited)
		}
	}
}

func clearParameterDescriptions(parameters openapi3.Parameters, visited map[*openapi3.Schema]bool) {
	for _, paramRef := range parameters {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value
		param.Description = ""
		clearSchemaDescriptions(param.Schema, visited)
		clearContentDescriptions(param.Content, visited)
	}
}

func clearExtensions(spec *openapi3.T) {
	visited := make(map[*openapi3.Schema]bool)

	spec.Extensions = nil
	if spec.Components != nil {
		spec.Components.Extensions = nil
		for _, ref := range spec.Components.Schemas {
			clearSchemaExtensions(ref, visited)
		}
	}

	for _, path := range spec.Paths.Map() {
		if path == nil {
			continue
		}

		path.Extensions = nil
		clearParameterExtensions(path.Parameters, visited)

		for _, operation := range path.Operations() {
			if operation == nil {
				continue
			}

			operation.Extensions = nil
			clearParameterExtensions(operation.Parameters, visited)

			if operation.RequestBody != nil {
				bodyRef := operation.RequestBody

				bodyRef.Extensions = nil
				if bodyRef.Value != nil {
					bodyRef.Value.Extensions = nil
					clearContentExtensions(bodyRef.Value.Content, visited)
				}
			}

			if operation.Responses != nil {
				operation.Responses.Extensions = nil
				for _, responseRef := range operation.Responses.Map() {
					if responseRef == nil {
						continue
					}

					responseRef.Extensions = nil
					if responseRef.Value == nil {
						continue
					}

					response := responseRef.Value

					response.Extensions = nil
					for _, headerRef := range response.Headers {
						if headerRef == nil {
							continue
						}

						headerRef.Extensions = nil
						if headerRef.Value != nil {
							header := headerRef.Value
							header.Extensions = nil
							clearSchemaExtensions(header.Schema, visited)
							clearContentExtensions(header.Content, visited)
						}
					}

					clearContentExtensions(response.Content, visited)
				}
			}
		}
	}
}

func clearSchemaExtensions(ref *openapi3.SchemaRef, visited map[*openapi3.Schema]bool) {
	if ref == nil {
		return
	}

	ref.Extensions = nil
	if ref.Value == nil || visited[ref.Value] {
		return
	}

	schema := ref.Value
	visited[schema] = true

	schema.Extensions = nil
	for _, property := range schema.Properties {
		clearSchemaExtensions(property, visited)
	}

	for _, refs := range [][]*openapi3.SchemaRef{
		schema.AllOf,
		schema.AnyOf,
		schema.OneOf,
	} {
		for _, nested := range refs {
			clearSchemaExtensions(nested, visited)
		}
	}

	clearSchemaExtensions(schema.Items, visited)
	clearSchemaExtensions(schema.Not, visited)
	clearSchemaExtensions(schema.AdditionalProperties.Schema, visited)
}

func clearContentExtensions(content openapi3.Content, visited map[*openapi3.Schema]bool) {
	for _, media := range content {
		if media == nil {
			continue
		}

		media.Extensions = nil
		clearSchemaExtensions(media.Schema, visited)
	}
}

func clearParameterExtensions(parameters openapi3.Parameters, visited map[*openapi3.Schema]bool) {
	for _, paramRef := range parameters {
		if paramRef == nil {
			continue
		}

		paramRef.Extensions = nil
		if paramRef.Value == nil {
			continue
		}

		param := paramRef.Value
		param.Extensions = nil
		clearSchemaExtensions(param.Schema, visited)
		clearContentExtensions(param.Content, visited)
	}
}

func clearExamples(spec *openapi3.T) {
	visited := make(map[*openapi3.Schema]bool)

	if spec.Components != nil {
		for _, ref := range spec.Components.Schemas {
			clearSchemaExamples(ref, visited)
		}

		for _, paramRef := range spec.Components.Parameters {
			clearParameterExample(paramRef, visited)
		}

		for _, headerRef := range spec.Components.Headers {
			clearHeaderExamples(headerRef, visited)
		}

		for _, bodyRef := range spec.Components.RequestBodies {
			clearRequestBodyExamples(bodyRef, visited)
		}

		for _, responseRef := range spec.Components.Responses {
			clearResponseExamples(responseRef, visited)
		}

		for _, exampleRef := range spec.Components.Examples {
			clearExampleRef(exampleRef)
		}
	}

	for _, path := range spec.Paths.Map() {
		if path == nil {
			continue
		}

		clearParameterExamples(path.Parameters, visited)

		for _, operation := range path.Operations() {
			if operation == nil {
				continue
			}

			clearParameterExamples(operation.Parameters, visited)
			clearRequestBodyExamples(operation.RequestBody, visited)

			if operation.Responses != nil {
				for _, responseRef := range operation.Responses.Map() {
					clearResponseExamples(responseRef, visited)
				}
			}
		}
	}
}

func clearSchemaExamples(ref *openapi3.SchemaRef, visited map[*openapi3.Schema]bool) {
	if ref == nil || ref.Value == nil || visited[ref.Value] {
		return
	}

	schema := ref.Value
	visited[schema] = true

	schema.Example = nil
	schema.Examples = nil

	for _, property := range schema.Properties {
		clearSchemaExamples(property, visited)
	}

	for _, refs := range [][]*openapi3.SchemaRef{
		schema.AllOf,
		schema.AnyOf,
		schema.OneOf,
		schema.PrefixItems,
	} {
		for _, nested := range refs {
			clearSchemaExamples(nested, visited)
		}
	}

	for _, ref := range schema.Defs {
		clearSchemaExamples(ref, visited)
	}

	clearSchemaExamples(schema.Items, visited)
	clearSchemaExamples(schema.Not, visited)
	clearSchemaExamples(schema.Contains, visited)
	clearSchemaExamples(schema.ContentSchema, visited)
	clearSchemaExamples(schema.AdditionalProperties.Schema, visited)
}

func clearExampleRef(ref *openapi3.ExampleRef) {
	if ref == nil || ref.Value == nil {
		return
	}

	ref.Value.Value = nil
	ref.Value.ExternalValue = ""
}

func clearContentExamples(content openapi3.Content, visited map[*openapi3.Schema]bool) {
	for _, media := range content {
		if media == nil {
			continue
		}

		media.Example = nil

		for _, exampleRef := range media.Examples {
			clearExampleRef(exampleRef)
		}

		media.Examples = nil

		clearSchemaExamples(media.Schema, visited)
	}
}

func clearParameterExample(paramRef *openapi3.ParameterRef, visited map[*openapi3.Schema]bool) {
	if paramRef == nil || paramRef.Value == nil {
		return
	}

	param := paramRef.Value
	param.Example = nil

	for _, exampleRef := range param.Examples {
		clearExampleRef(exampleRef)
	}

	param.Examples = nil

	clearSchemaExamples(param.Schema, visited)
	clearContentExamples(param.Content, visited)
}

func clearParameterExamples(parameters openapi3.Parameters, visited map[*openapi3.Schema]bool) {
	for _, paramRef := range parameters {
		clearParameterExample(paramRef, visited)
	}
}

func clearHeaderExamples(headerRef *openapi3.HeaderRef, visited map[*openapi3.Schema]bool) {
	if headerRef == nil || headerRef.Value == nil {
		return
	}

	header := headerRef.Value
	header.Example = nil

	for _, exampleRef := range header.Examples {
		clearExampleRef(exampleRef)
	}

	header.Examples = nil

	clearSchemaExamples(header.Schema, visited)
	clearContentExamples(header.Content, visited)
}

func clearRequestBodyExamples(bodyRef *openapi3.RequestBodyRef, visited map[*openapi3.Schema]bool) {
	if bodyRef == nil || bodyRef.Value == nil {
		return
	}

	clearContentExamples(bodyRef.Value.Content, visited)
}

func clearResponseExamples(responseRef *openapi3.ResponseRef, visited map[*openapi3.Schema]bool) {
	if responseRef == nil || responseRef.Value == nil {
		return
	}

	response := responseRef.Value

	for _, headerRef := range response.Headers {
		clearHeaderExamples(headerRef, visited)
	}

	clearContentExamples(response.Content, visited)
}
