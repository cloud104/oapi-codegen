// internal/naming/inflection_test.go
package naming

import "testing"

func TestSingular(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Regular plurals.
		{"users", "user"},
		{"books", "book"},
		{"categories", "category"},
		{"boxes", "box"},
		{"statuses", "status"},

		// Irregular plurals.
		{"people", "person"},
		{"children", "child"},
		{"men", "man"},
		{"mice", "mouse"},

		// Already singular.
		{"user", "user"},
		{"person", "person"},

		// Uncountable nouns.
		{"sheep", "sheep"},
		{"fish", "fish"},
		{"equipment", "equipment"},

		// Initialisms.
		{"TKS", "TKS"},
		{"AWS", "AWS"},
		{"DNS", "DNS"},
		{"API", "API"},
		{"URL", "URL"},

		// Plural initialisms.
		{"APIs", "API"},
		{"URLs", "URL"},

		// PascalCase.
		{"Users", "User"},
		{"Categories", "Category"},

		// Empty input.
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Singular(tt.input); got != tt.want {
				t.Errorf("Singular(%q) = %q; want %q",
					tt.input, got, tt.want)
			}
		})
	}
}
