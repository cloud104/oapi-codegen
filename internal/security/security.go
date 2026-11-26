package security

import (
	"strings"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

type Security struct {
	APIKey []*APIKeyConfig
	Basic  []*BasicConfig
	Bearer []*BearerConfig
	OAuth2 []*OAuth2Config
}

func ForOperation(op codegen.OperationDefinition) *Security {
	return ForOperations([]codegen.OperationDefinition{op})
}

func ForOperations(ops []codegen.OperationDefinition) *Security {
	s := &Security{}

	for _, auth := range schemesForOperations(ops) {
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

		case "oauth2":
			s.addOAuth2(auth)
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

func (s *Security) HasOAuth2() bool {
	return len(s.OAuth2) > 0
}

func (s *Security) HasSecurity() bool {
	return s.HasAPIKey() ||
		s.HasBasic() ||
		s.HasBearer() ||
		s.HasOAuth2()
}
