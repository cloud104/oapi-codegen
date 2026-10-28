package codegenstate

import "github.com/getkin/kin-openapi/openapi3"

func Components() (*openapi3.Components, bool) {
	if globalState.spec == nil || globalState.spec.Components == nil {
		return nil, false
	}

	return globalState.spec.Components, true
}
