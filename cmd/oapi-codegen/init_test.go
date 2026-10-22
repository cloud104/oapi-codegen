package main

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"math"
	"encoding/json"
)

func Test_getMaxBodySize(t *testing.T) {
	tests := []struct {
		name string
		size any
		want int64
	}{
		// String values
		{"valid MiB", "5 MiB", 5 * 1024 * 1024},
		{"valid KiB", "5 KiB", 5 * 1024},
		{"valid GiB", "2 GiB", 2 * 1024 * 1024 * 1024},
		{"bytes", "1024 B", 1024},
		{"numeric string", "524288", 524288},
		{"size without space", "5MiB", 5 * 1024 * 1024},
		{"size with whitespace", " 5 MiB ", 5 * 1024 * 1024},
		{"minimum size", "1 B", 1},
		{"maximum size", "9223372036854775807 B", math.MaxInt64},

		// Invalid string values
		{"empty string", "", defaultMaxBodySize},
		{"whitespace only", "   ", defaultMaxBodySize},
		{"invalid string", "invalid", defaultMaxBodySize},
		{"negative string", "-1 B", defaultMaxBodySize},
		{"zero string", "0 B", defaultMaxBodySize},
		{"exceeds maximum", "9223372036854775808 B", defaultMaxBodySize},

		// Signed integers
		{"int", int(524288), 524288},
		{"int8", int8(127), 127},
		{"int16", int16(32767), 32767},
		{"int32", int32(524288), 524288},
		{"int64", int64(524288), 524288},
		{"minimum integer", int64(1), 1},
		{"maximum integer", int64(math.MaxInt64), math.MaxInt64},
		{"zero integer", int(0), defaultMaxBodySize},
		{"negative integer", int(-1), defaultMaxBodySize},
		{"minimum int64", int64(math.MinInt64), defaultMaxBodySize},

		// Unsigned integers
		{"uint", uint(524288), 524288},
		{"uint8", uint8(255), 255},
		{"uint16", uint16(65535), 65535},
		{"uint32", uint32(524288), 524288},
		{"uint64", uint64(524288), 524288},
		{"maximum uint64 within range", uint64(math.MaxInt64), math.MaxInt64},
		{"uint64 exceeds maximum", uint64(math.MaxInt64) + 1, defaultMaxBodySize},
		{"maximum uint64", uint64(math.MaxUint64), defaultMaxBodySize},
		{"zero unsigned integer", uint(0), defaultMaxBodySize},

		// Floating-point values
		{"float32", float32(1024), 1024},
		{"float64", float64(524288), 524288},
		{"minimum float", float64(1), 1},
		{"fractional float32", float32(1.5), 1},
		{"fractional float64", float64(1024.5), 1024},
		{"zero float", float64(0), defaultMaxBodySize},
		{"negative float", float64(-1), defaultMaxBodySize},
		{"float exceeds maximum", float64(0x1p63), defaultMaxBodySize},
		{"positive infinity", math.Inf(1), defaultMaxBodySize},
		{"negative infinity", math.Inf(-1), defaultMaxBodySize},
		{"NaN", math.NaN(), defaultMaxBodySize},

		// JSON numbers
		{"JSON integer", json.Number("524288"), 524288},
		{"JSON minimum", json.Number("1"), 1},
		{"JSON maximum", json.Number("9223372036854775807"), math.MaxInt64},
		{"JSON zero", json.Number("0"), defaultMaxBodySize},
		{"JSON negative", json.Number("-1"), defaultMaxBodySize},
		{"JSON exceeds maximum", json.Number("9223372036854775808"), defaultMaxBodySize},
		{"JSON fractional", json.Number("1.5"), 1},
		{"JSON invalid", json.Number("invalid"), defaultMaxBodySize},
		{"JSON empty", json.Number(""), defaultMaxBodySize},

		// Missing or unsupported values
		{"missing extension", nil, defaultMaxBodySize},
		{"boolean true", true, defaultMaxBodySize},
		{"boolean false", false, defaultMaxBodySize},
		{"slice", []int{1024}, defaultMaxBodySize},
		{"map", map[string]any{"size": 1024}, defaultMaxBodySize},
		{"struct", struct{}{}, defaultMaxBodySize},
		{"pointer", new(int64), defaultMaxBodySize},
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

			if got := bodyLimit(op).Bytes(); got != tt.want {
				t.Errorf("getMaxBodySize() = %d, want %d", got, tt.want)
			}
		})
	}
}
