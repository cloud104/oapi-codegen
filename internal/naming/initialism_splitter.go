package naming

import (
	"strings"
	"unicode"
	_ "unsafe"

	"github.com/ettle/strcase"
)

//go:linkname defaultSplitFn github.com/ettle/strcase.defaultSplitFn
func defaultSplitFn(prev, curr, next rune) strcase.SplitAction

// initialismSplitter implements a strcase.SplitFn that detects initialism boundaries.
type initialismSplitter struct {
	initialisms map[string]bool
	word        string
}

func newInitialismSplitter(initialisms map[string]bool) *initialismSplitter {
	return &initialismSplitter{initialisms: initialisms}
}

func (s *initialismSplitter) split(prev, curr, next rune) strcase.SplitAction {
	// Reset at the start of input.
	if prev == 0 {
		s.word = ""
	}

	// Check for an initialism boundary.
	action := defaultSplitFn(prev, curr, next)
	if action == strcase.Noop && s.isBoundary(prev, curr) {
		action = strcase.Split
	}

	// Update the current word.
	switch action {
	case strcase.Split:
		s.word = strings.ToUpper(string(curr))
	case strcase.SkipSplit:
		s.word = ""
	case strcase.Skip:
	default:
		s.word += strings.ToUpper(string(curr))
	}

	return action
}

func (s *initialismSplitter) isBoundary(prev, curr rune) bool {
	// Require uppercase letters and a known initialism.
	if !unicode.IsUpper(prev) || !unicode.IsUpper(curr) || !s.initialisms[s.word] {
		return false
	}

	// Check for longer matching initialisms.
	prefix := s.word + string(curr)
	for name := range s.initialisms {
		if strings.HasPrefix(name, prefix) {
			return false
		}
	}

	return true
}
