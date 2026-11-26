package security

type BearerConfig struct {
	Name         string
	Description  string
	BearerFormat string
	Scopes       []string
}

func (s *Security) addBearer(auth *OperationScheme) {
	s.Bearer = append(s.Bearer, &BearerConfig{
		Name:         auth.Definition.ProviderName,
		Description:  auth.Scheme.Description,
		BearerFormat: auth.Scheme.BearerFormat,
		Scopes:       auth.Definition.Scopes,
	})
}
