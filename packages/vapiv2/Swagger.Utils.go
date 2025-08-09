package vapi

type SwaggerUtilsType struct {
	Oauth2 map[string]interface{}
}

func (s *SwaggerUtilsType) Oauth2AuthCodePKCE(
	authorizationUrl string,
	tokenUrl string,
	scopes map[string]string,
	description string,
) map[string]interface{} {
	if description == "" {
		description = "\"OAuth2 Authorization Code Flow with PKCE support\""
	}
	oauth2AuthCodePKCE := map[string]interface{}{
		"type":             "oauth2",
		"flow":             "accessCode",
		"authorizationUrl": authorizationUrl,
		"tokenUrl":         tokenUrl,
		"scopes":           scopes,
		"description":      description,
	}
	s.Oauth2["OAuth2AuthCodePKCE"] = oauth2AuthCodePKCE
	return oauth2AuthCodePKCE
}

var SwaggerUtils = &SwaggerUtilsType{
	Oauth2: map[string]interface{}{},
}

func (s *SwaggerUtilsType) OAuth2Password(tokenUrl string, description string) map[string]interface{} {
	if description == "" {
		description = "\"OAuth2 Password Flow - Enter email/username and password in the popup to get token.\""
	}
	oauth2Password := map[string]interface{}{
		"description": description,
		"flow":        "password",
		"tokenUrl":    tokenUrl,
		"type":        "oauth2",
	}
	s.Oauth2["OAuth2Password"] = oauth2Password
	return oauth2Password

}
