package codegenstate

import "github.com/getkin/kin-openapi/openapi3"

func Security() (openapi3.SecurityRequirements, bool) {
	if globalState.spec == nil || globalState.spec.Security == nil {
		return nil, false
	}

	return globalState.spec.Security, true
}
