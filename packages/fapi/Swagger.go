package fapi

import (
	"embed"
	"encoding/json"
	"fapi/swaggers"
	"sync"
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
	swaggerData.BasePath = basrUrl
	LoadHandlerInfo(swaggerData)
	ret, err := json.Marshal(swaggerData)

	if err != nil {
		panic(err)
	}

	return ret
}
