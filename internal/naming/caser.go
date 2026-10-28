package naming

import (
	"strings"

	"github.com/ettle/strcase"
)

type Caser struct {
	converter *strcase.Caser
}

func NewCaser(initialisms []string) *Caser {
	initialismOverrides := make(map[string]bool, len(initialisms))
	for _, word := range initialisms {
		initialismOverrides[strings.ToUpper(word)] = true
	}

	splitter := newInitialismSplitter(initialismOverrides)
	caser := strcase.NewCaser(true, initialismOverrides, splitter.split)

	return &Caser{converter: caser}
}

func (c *Caser) Camel(name string) string {
	return c.converter.ToCamel(name)
}

func (c *Caser) Pascal(name string) string {
	return c.converter.ToPascal(name)
}
