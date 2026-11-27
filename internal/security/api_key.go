package security

type APIKey struct {
	Name        string
	Description string
	KeyName     string
	In          string
	Scopes      []string
}

func (s *Security) addAPIKey(auth *OperationScheme) {
	s.APIKey = append(s.APIKey, &APIKey{
		Name:        auth.Definition.ProviderName,
		Description: auth.Scheme.Description,
		KeyName:     auth.Scheme.Name,
		In:          auth.Scheme.In,
		Scopes:      auth.Definition.Scopes,
	})
}
