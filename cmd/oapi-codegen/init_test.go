package main

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func Test_getMaxBodySize(t *testing.T) {
	tests := []struct {
		name string
		size any
		want int64
	}{
		{"valid size", "5 MiB", 5 * 1024 * 1024},
		{"size in bytes", "1024 B", 1024},
		{"missing extension", nil, defaultMaxBodySize},
		{"invalid size", "invalid", defaultMaxBodySize},
		{"zero size", "0 B", defaultMaxBodySize},
		{"exceeds maximum", "9223372036854775808 B", defaultMaxBodySize},
		{"", "524288", 524288},
		{"", 524288, 524288},
		{"", "5MiB", 5242880},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &openapi3.RequestBody{}
			if tt.size != nil {
				body.Extensions = map[string]any{
					maxBodySizeExtension: tt.size,
				}
			}

			op := codegen.OperationDefinition{
				Spec: &openapi3.Operation{
					RequestBody: &openapi3.RequestBodyRef{Value: body},
				},
			}

			if got := getMaxBodySize(op).Bytes(); got != tt.want {
				t.Errorf("getMaxBodySize() = %d, want %d", got, tt.want)
			}
		})
	}
}
