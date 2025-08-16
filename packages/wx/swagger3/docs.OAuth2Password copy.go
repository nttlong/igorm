package swaggers

func (s *Swagger) OAuth2Password(tokenUrl string) *Swagger {
	s.Components.SecuritySchemes["OAuth2Password"] = SecurityScheme{
		Type: "oauth2",
		Flows: &OAuthFlows{
			Password: &OAuthFlow{
				TokenUrl: tokenUrl,
			},
		},
	}
	return s

}
