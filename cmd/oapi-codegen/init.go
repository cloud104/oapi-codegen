package main

import (
	sprig "github.com/Masterminds/sprig/v3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

func init() {
	for name, fn := range sprig.FuncMap() {
		if _, exists := codegen.TemplateFunctions[name]; !exists {
			codegen.TemplateFunctions[name] = fn
		}
	}
}
