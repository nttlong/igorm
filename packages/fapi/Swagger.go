package fapi

import (
	"embed"
	"encoding/json"
	"sync"
)

// Parameter định nghĩa một tham số của request
type Parameter struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
}

// AdditionalProperties được dùng cho các object có key động
type AdditionalProperties struct {
	Type string `json:"type"`
}

// Schema là cấu trúc dữ liệu của một đối tượng JSON
// Nó có thể là một định nghĩa inline hoặc một reference ($ref)
type Schema struct {
	// Dành cho trường hợp `$ref`
	Ref *string `json:"$ref,omitempty"`

	// Dành cho trường hợp định nghĩa schema inline
	Type                 string                `json:"type,omitempty"`
	Properties           map[string]Schema     `json:"properties,omitempty"`
	AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty"`
}

// Response định nghĩa cấu trúc của một response HTTP
type Response struct {
	Description string  `json:"description"`
	Schema      *Schema `json:"schema,omitempty"`
}
type Operation struct {
	Consumes   []string              `json:"consumes"`
	Produces   []string              `json:"produces"`
	Tags       []string              `json:"tags"`
	Summary    string                `json:"summary"`
	Parameters []Parameter           `json:"parameters"`
	Responses  map[string]Response   `json:"responses"`
	Security   []map[string][]string `json:"security"`
}

// PathItem đại diện cho một đường dẫn API (ví dụ: "/api/auth/login")
// Key của map là phương thức HTTP (ví dụ: "post", "get")
type PathItem struct {
	Post    *Operation `json:"post,omitempty"`
	Get     *Operation `json:"get,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	// ...
	// Có thể thêm các phương thức khác như Get, Put, Delete...
}

// Definition định nghĩa các cấu trúc dữ liệu có thể tái sử dụng
// Ví dụ: "controllers.OAuthTokenResponse"
type Definition struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

// Property là một trường trong Definition
type Property struct {
	Type string `json:"type"`
}

// SecurityDefinition chứa các cấu hình bảo mật
type SecurityDefinition struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Flow        string `json:"flow"`
	TokenURL    string `json:"tokenUrl"`
}

// Swagger là cấu trúc chính để biểu diễn toàn bộ file Swagger 2.0
type Swagger struct {
	Swagger             string                        `json:"swagger"`
	Info                Info                          `json:"info"`
	Host                string                        `json:"host"`
	BasePath            string                        `json:"basePath"`
	Paths               map[string]PathItem           `json:"paths"`
	Definitions         map[string]Definition         `json:"definitions"`
	SecurityDefinitions map[string]SecurityDefinition `json:"securityDefinitions"`
}

// Info chứa các thông tin về API
type Info struct {
	Description string   `json:"description"`
	Title       string   `json:"title"`
	Contact     struct{} `json:"contact"`
	Version     string   `json:"version"`
}

// createMockSwaggerJSON là một hàm helper để tạo file swagger.json cho ví dụ này.
// Trong thực tế, bạn sẽ có file swagger.json được tạo ra bởi các công cụ khác.
func CreateMockSwaggerJSON() []byte {
	swaggerData.Paths["api/hello"] = PathItem{
		Post: &Operation{
			Consumes:   []string{"application/x-www-form-urlencoded", "application/json"},
			Produces:   []string{"application/x-www-form-urlencoded", "application/json"},
			Parameters: []Parameter{},
			Responses:  map[string]Response{},
		},
	}

	ret, err := json.Marshal(swaggerData)

	if err != nil {
		panic(err)
	}

	return ret
}

//go:embed swagger.json
var embeddedSwaggerJSON embed.FS
var loadSwaggerInfoOnce sync.Once

func loadSwaggerInfo() (*Swagger, error) {

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

var swaggerData *Swagger

func init() {
	swaggerData, _ = loadSwaggerInfo()
}
