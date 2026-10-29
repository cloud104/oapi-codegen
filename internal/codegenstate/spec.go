package codegenstate

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pb33f/libopenapi/bundler"
	"github.com/pb33f/libopenapi/datamodel"
)

var (
	inlineSpecOnce sync.Once
	inlineSpec     *openapi3.T
	inlineSpecErr  error
)

func InlineSpec() (*openapi3.T, error) {
	inlineSpecOnce.Do(func() {
		data, err := json.Marshal(globalState.spec)
		if err != nil {
			inlineSpecErr = fmt.Errorf("marshal OpenAPI spec: %w", err)
			return
		}

		bundled, err := bundler.BundleBytes(data, &datamodel.DocumentConfiguration{
			AllowFileReferences: true,
			BundleInlineRefs:    true,
		})
		if err != nil {
			inlineSpecErr = fmt.Errorf("bundle OpenAPI spec: %w", err)
			return
		}

		inlineSpec, inlineSpecErr = openapi3.NewLoader().LoadFromData(bundled)
		if inlineSpecErr != nil {
			inlineSpecErr = fmt.Errorf("load bundled OpenAPI spec: %w", inlineSpecErr)
		}
	})

	return inlineSpec, inlineSpecErr
}
