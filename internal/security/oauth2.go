package security

type OAuth2 struct {
	Name        string
	Description string
	Flow        string
	Scopes      []string
	TokenURL    string
}

func (s *Security) addOAuth2(auth *OperationScheme) {
	flows := auth.Scheme.Flows
	if flows == nil {
		return
	}

	if flow := flows.ClientCredentials; flow != nil {
		s.OAuth2 = append(s.OAuth2, &OAuth2{
			Name:        auth.Definition.ProviderName,
			Description: auth.Scheme.Description,
			Flow:        "clientCredentials",
			Scopes:      auth.Definition.Scopes,
			TokenURL:    flow.TokenURL,
		})
	}

	if flow := flows.Password; flow != nil {
		s.OAuth2 = append(s.OAuth2, &OAuth2{
			Name:        auth.Definition.ProviderName,
			Description: auth.Scheme.Description,
			Flow:        "password",
			Scopes:      auth.Definition.Scopes,
			TokenURL:    flow.TokenURL,
		})
	}
}
