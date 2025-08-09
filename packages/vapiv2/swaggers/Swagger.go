package swaggers

// Property là một trường trong Definition
type Property struct {
	Type string `json:"type"`
}

// Definition định nghĩa các cấu trúc dữ liệu có thể tái sử dụng
// Ví dụ: "controllers.OAuthTokenResponse"
type Definition struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
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
	Swagger     string                `json:"swagger"`
	Info        Info                  `json:"info"`
	Host        string                `json:"host"`
	BasePath    string                `json:"basePath"`
	Paths       map[string]PathItem   `json:"paths"`
	Definitions map[string]Definition `json:"definitions"`
	// SecurityDefinitions map[string]SecurityDefinition `json:"securityDefinitions"`
	SecurityDefinitions map[string]interface{} `json:"securityDefinitions"`
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
