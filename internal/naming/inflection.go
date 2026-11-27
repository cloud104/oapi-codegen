package naming

import "github.com/jinzhu/inflection"

func Singular(str string) string {
	return inflection.Singular(str)
}
