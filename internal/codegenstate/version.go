package codegenstate

import (
	"errors"
	"fmt"

	semver "github.com/Masterminds/semver/v3"
)

var openapi31 = semver.MustParse("3.1.0")

// Version returns the OpenAPI specification version as a semantic version.
func Version() (*semver.Version, error) {
	spec := globalState.spec
	if spec == nil {
		return nil, errors.New("OpenAPI specification is not initialized")
	}

	version, err := semver.StrictNewVersion(spec.OpenAPI)
	if err != nil {
		return nil, fmt.Errorf("invalid OpenAPI version %q: %w", spec.OpenAPI, err)
	}

	return version, nil
}

// IsJSONSchemaCompatible reports whether the OpenAPI version uses
// a Schema Object compatible with standard JSON Schema.
//
// OpenAPI 3.1 and later align their Schema Object with JSON Schema,
// whereas earlier versions use an extended subset with restricted
// keywords, modified semantics, and OpenAPI-specific extensions.
func IsJSONSchemaCompatible() bool {
	v, err := Version()
	if err != nil {
		panic(err)
	}

	return v.GreaterThanEqual(openapi31)
}
