package naming

import (
	"slices"
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/cloud104/oapi-codegen/internal/codegenstate"
)

func Singular(str string) (result string) {
	// Preserve registered initialisms and empty strings.
	if str == "" || isInitialism(str) {
		return str
	}

	// Remove the plural suffix from initialisms, such as APIs and URLs.
	if isPluralInitialism(str) {
		return strings.TrimSuffix(str, "s")
	}

	return inflection.Singular(str)
}

func Plural(str string) (result string) {
	// Preserve registered initialisms and empty strings.
	if str == "" || isInitialism(str) {
		return str
	}

	// Preserve already pluralized initialisms, such as APIs and URLs.
	if isPluralInitialism(str) {
		return str
	}

	// Normalize irregular plurals before pluralizing to prevent
	// incorrect results such as "peoples" and "childrens".
	return inflection.Plural(inflection.Singular(str))
}

func isPluralInitialism(str string) bool {
	return strings.HasSuffix(str, "s") &&
		isInitialism(strings.TrimSuffix(str, "s"))
}

func isInitialism(str string) bool {
	initialisms := codegenstate.Initialisms()
	return slices.Contains(initialisms, strings.ToLower(str))
}
