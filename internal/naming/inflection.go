package naming

import (
	"strings"

	"github.com/jinzhu/inflection"
)

func Singular(str string) string {
	// Preserve uppercase initialisms and empty strings.
	if str == strings.ToUpper(str) {
		return str
	}

	return inflection.Singular(str)
}

func Plural(str string) string {
	// Preserve uppercase initialisms and empty strings.
	if str == strings.ToUpper(str) {
		return str
	}

	singular := Singular(str)

	// Preserve already pluralized initialisms, such as APIs and URLs.
	if singular == strings.ToUpper(singular) {
		return str
	}

	// Normalize irregular plurals before pluralizing to prevent
	// incorrect results such as "peoples" and "childrens".
	return inflection.Plural(singular)
}
