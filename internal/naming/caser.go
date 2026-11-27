package naming

import (
	"strings"

	"github.com/ettle/strcase"

	"github.com/cloud104/oapi-codegen/internal/codegenstate"
)

func Camel(name string) string {
	caser := newCaser()
	return caser.ToCamel(name)
}

func Pascal(name string) string {
	caser := newCaser()
	return caser.ToPascal(name)
}

func newCaser() *strcase.Caser {
	initialisms := codegenstate.Initialisms()
	overrides := make(map[string]bool, len(initialisms))

	for _, word := range initialisms {
		overrides[strings.ToUpper(word)] = true
	}

	return strcase.NewCaser(true, overrides, newInitialismSplitter(overrides).split)
}
