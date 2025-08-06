package swaggers

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
