package swaggers

// Parameter định nghĩa một tham số của request
type Parameter struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Name        string `json:"name"`
	/*
					  Swagger input control
					  Value is 'body',"formBody" or "path"
					  path	Dữ liệu nằm trong URL path	/file/{FileID}
		query	Dữ liệu nằm trong query string	/file?FileID=123
		header	Dữ liệu được truyền qua HTTP header	X-Request-ID: abc123
		cookie	Dữ liệu nằm trong cookie	Cookie: FileID=xyz
		body	in: body	➜ dùng requestBody
		formData	in: formData	➜ dùng requestBody với content-type application/x-www-form-urlencoded hoặc multipart/form-data
		"path" if it in url
	*/
	In       string `json:"in"`
	Required bool   `json:"required"`
	Example  string `json:"example"`
}
