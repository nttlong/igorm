package htttpserver

import (
	_ "embed"
	"encoding/json"
	"net/http"
	swaggers3 "wx/swagger3"
	//httpSwagger "github.com/swaggo/http-swagger"
)

//go:embed swagger/index.html
var indexHtml []byte

//go:embed swagger/swagger-ui.css
var css []byte

//go:embed swagger/swagger-ui-bundle.js
var js []byte

type SwaggerContact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}
type SwaggerInfo struct {
	Title       string          `json:"title"`
	Description string          `json:"description,omitempty"`
	Version     string          `json:"version"`
	Contact     *SwaggerContact `json:"contact,omitempty"`
}

/*
 */
type OAuth2AuthCodePKCE struct {
}

func CreateSwagger(server *HtttpServer, BaseUri string) SwaggerBuild {
	sw, err := swaggers3.CreateSwagger(server.BaseUrl, swaggers3.Info{})
	if err != nil {
		return SwaggerBuild{
			server:  server,
			BaseUri: BaseUri,
			err:     err,
		}
	}
	return SwaggerBuild{
		server:  server,
		BaseUri: BaseUri,
		swagger: sw,
	}
}
func (sb *SwaggerBuild) Info(info SwaggerInfo) *SwaggerBuild {
	sb.info = info
	return sb
}
func (sb *SwaggerBuild) OAuth2AuthCodePKCE(AuthorizationUrl string, TokenUrl string, Scopes map[string]string) *SwaggerBuild {
	sb.swagger.OAuth2AuthCodePKCE(AuthorizationUrl, TokenUrl, Scopes)
	return sb
}

/*
Enable OAuth2 Password flow on Swagger docs.

@param TokenUrl the URL to obtain the token
*/
func (sb *SwaggerBuild) OAuth2Password(TokenUrl string) *SwaggerBuild {
	sb.swagger.OAuth2Password(TokenUrl)
	return sb
}
func (sb *SwaggerBuild) Build() error {
	server := sb.server
	useSwagger = true
	mux := server.mux
	uri := sb.BaseUri
	//sb.swagger3GetPaths()
	sb.LoadFromRoutes()
	data, err := json.Marshal(sb.swagger)
	if err != nil {
		sb.err = err

	}

	// info := sb.info
	// 1. Phục vụ file swagger.json từ đường dẫn /swagger.json

	mux.HandleFunc(uri+"/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write(indexHtml); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})

	/*

		Serve the swagger.json file from the path /swagger.json
		// The httpSwagger library will look for this file to display the documentation.
	*/
	mux.HandleFunc(uri+"/swagger.json", func(w http.ResponseWriter, r *http.Request) {

		if sb.err != nil {
			http.Error(w, sb.err.Error(), http.StatusInternalServerError)
			return
		}

		// Thiết lập header để trình duyệt hiểu đây là file JSON
		w.Header().Set("Content-Type", "application/json")

		if _, err := w.Write(data); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})
	/*

		Serve the swagger-ui.css file from the path swagger-ui.css
		// The httpSwagger library will look for this file to display the documentation.
	*/
	mux.HandleFunc(uri+"/swagger-ui.css", func(w http.ResponseWriter, r *http.Request) {
		// Đọc file swagger.json từ thư mục hiện tại
		w.Header().Set("Content-Type", "text/css")

		if _, err := w.Write(css); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc(uri+"/swagger-ui-bundle.js", func(w http.ResponseWriter, r *http.Request) {
		// Đọc file swagger.json từ thư mục hiện tại
		w.Header().Set("Content-Type", "application/javascript")
		if _, err := w.Write(js); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})
	// 2. Phục vụ giao diện Swagger UI trên đường dẫn /swagger/
	// Thư viện httpSwagger.WrapHandler tự động tạo giao diện HTML.
	// Đường dẫn thứ hai "./swagger.json" là vị trí của file JSON mà UI sẽ hiển thị.
	//mux.Handle("/swagger/", httpSwagger.WrapHandler)
	return nil
}
