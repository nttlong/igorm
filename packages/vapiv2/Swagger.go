package vapi

import (
	"embed"
	"encoding/json"
	"sync"
	"vapi/swaggers"
)

// Response định nghĩa cấu trúc của một response HTTP

//go:embed swagger.json
var embeddedSwaggerJSON embed.FS
var loadSwaggerInfoOnce sync.Once

func loadSwaggerInfo() (*swaggers.Swagger, error) {

	var err error
	loadSwaggerInfoOnce.Do(func() {

		data, err := embeddedSwaggerJSON.ReadFile("swagger.json")
		if err != nil {
			return
		}

		err = json.Unmarshal(data, &swaggerData)
		if err != nil {
			return
		}

	})
	return swaggerData, err
}

var swaggerData *swaggers.Swagger

func init() {
	swaggerData, _ = loadSwaggerInfo()
	//
	oauth2AuthCodePKCE := map[string]interface{}{
		"type":             "oauth2",
		"flow":             "accessCode",
		"authorizationUrl": "/oauth/authorize",
		"tokenUrl":         "/oauth/token",
		"scopes": map[string]interface{}{
			"read":  "Read access",
			"write": "Write access",
		},
		"description": "OAuth2 Authorization Code Flow with PKCE support",
	}
	swaggerData.SecurityDefinitions["OAuth2AuthCodePKCE"] = oauth2AuthCodePKCE
}

// createMockSwaggerJSON là một hàm helper để tạo file swagger.json cho ví dụ này.
// Trong thực tế, bạn sẽ có file swagger.json được tạo ra bởi các công cụ khác.
func CreateMockSwaggerJSON(basrUrl string) []byte {
	// swaggerData.Paths["api/hello"] = PathItem{
	// 	Post: &Operation{
	// 		Consumes:   []string{"application/x-www-form-urlencoded", "application/json"},
	// 		Produces:   []string{"application/x-www-form-urlencoded", "application/json"},
	// 		Parameters: []Parameter{},
	// 		Responses:  map[string]Response{},
	// 	},
	// }
	key := "CreateMockSwaggerJSON/" + basrUrl
	ret, _ := OnceCall(key, func() (*[]byte, error) {

		swaggerData.BasePath = basrUrl
		LoadHandlerInfo(swaggerData)
		for k, v := range SwaggerUtils.Oauth2 {
			swaggerData.SecurityDefinitions[k] = v
		}
		ret, err := json.Marshal(swaggerData)

		if err != nil {
			return nil, err
		}

		return &ret, nil
	})
	return *ret
}
