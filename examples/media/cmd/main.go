package main

import (
	"encoding/json"
	"fmt"
	"log"
	apps "media/internal/services/app"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"wx"
)

var rootPath = `D:\code\go\news2\igorm\examples\media\cmd\uploads`

func ListFiles(w http.ResponseWriter, r *http.Request) {
	// 1. Kiểm tra thư mục gốc
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		http.Error(w, "Directory not found.", http.StatusNotFound)
		return
	}

	// 2. Base URL
	baseUrl := "http://" + r.Host + r.URL.Path
	results := make([]string, 0, 1024) // pre-allocate

	// 3. WalkDir thay vì Walk
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == rootPath {
			return nil
		}

		// relative path
		rel, _ := filepath.Rel(rootPath, path)
		urlPath := strings.ReplaceAll(rel, "\\", "/")

		results = append(results, baseUrl+"/"+urlPath)
		return nil
	})

	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	// 4. Trả JSON
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false) // tránh encode `/` thành `\/`
	_ = enc.Encode(results)
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func main1() {
	http.HandleFunc("/api/media/hello", helloHandler)
	http.HandleFunc("/api/media/list-files", ListFiles)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
func main() {

	err := wx.Start(func(app *apps.App) error {
		go func() {
			f, _ := os.Create("mem.pprof")
			pprof.WriteHeapProfile(f)
			f.Close()
			log.Println("pprof listening on :6060")
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
		server, err := app.Server.Ins()
		if err != nil {
			return err
		}
		server.Start()
		return nil
	})
	if err != nil {
		panic(err)
	}
}
