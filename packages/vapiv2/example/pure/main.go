// file: main_http.go
package main

import (
	"net/http"
	"path/filepath"
)

func main() {
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		filename := "tes-004.pdf" // lấy tên file
		filePath := filepath.Join("./uploads", filename)
		http.ServeFile(w, r, filePath)
	})

	http.ListenAndServe(":8082", nil)
}
