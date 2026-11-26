package security

type BasicConfig struct {
	Name        string
	Description string
	Scopes      []string
}

func (s *Security) addBasic(auth *OperationScheme) {
	s.Basic = append(s.Basic, &BasicConfig{
		Name:        auth.Definition.ProviderName,
		Description: auth.Scheme.Description,
		Scopes:      auth.Definition.Scopes,
	})
}
