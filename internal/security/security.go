package security

import (
	"strings"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

type Security struct {
	APIKey []*APIKey
	Basic  []*Basic
	Bearer []*Bearer
}

func ForOperation(op codegen.OperationDefinition) *Security {
	return ForOperations([]codegen.OperationDefinition{op})
}

func ForOperations(ops []codegen.OperationDefinition) *Security {
	s := &Security{}

	for _, auth := range SchemesForOperations(ops) {
		switch auth.Scheme.Type {
		case "apiKey":
			s.addAPIKey(auth)

		case "http":
			switch strings.ToLower(auth.Scheme.Scheme) {
			case "basic":
				s.addBasic(auth)
			case "bearer":
				s.addBearer(auth)
			}
		}
	}

	return s
}

func (s *Security) HasAPIKey() bool {
	return len(s.APIKey) > 0
}

func (s *Security) HasBasic() bool {
	return len(s.Basic) > 0
}

func (s *Security) HasBearer() bool {
	return len(s.Bearer) > 0
}

func (s *Security) HasSecurity() bool {
	return s.HasAPIKey() ||
		s.HasBasic() ||
		s.HasBearer()
}
