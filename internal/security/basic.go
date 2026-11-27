package security

type Basic struct {
	Name        string
	Description string
	Scopes      []string
}

func (s *Security) addBasic(auth *OperationScheme) {
	s.Basic = append(s.Basic, &Basic{
		Name:        auth.Definition.ProviderName,
		Description: auth.Scheme.Description,
		Scopes:      auth.Definition.Scopes,
	})
}
