package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type Address struct {
	City     string `json:"city"`
	District string `json:"district"`
	Street   string `json:"street"`
}

type UserInput struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Email   string   `json:"email"`
	Phones  []string `json:"phones"`
	Address Address  `json:"address"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var input UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// giả lập xử lý
	msg := fmt.Sprintf("User %s, %d tuổi, sống ở %s",
		input.Name, input.Age, input.Address.City)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"msg": msg})
}
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Giới hạn dung lượng tối đa 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	// Lấy file từ form field "file"
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Lưu file lên ổ cứng tạm
	dst, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot create file: %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot save file: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Uploaded File: %s, Size: %d bytes, MIME: %s",
		handler.Filename, handler.Size, handler.Header.Get("Content-Type"))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File %s uploaded successfully!", handler.Filename)))
}

type TestApi struct{}

func (t *TestApi) DoSayHelloInVn(name string) string {
	return "xin chào, " + name
}
func (t *TestApi) DoSayHelloInEn(name string) string {
	return "hello " + name
}

func (t *TestApi) Hello(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	lang := chi.URLParam(r, "langCode")

	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if lang == "" {
		http.Error(w, "langCode is required", http.StatusBadRequest)
		return
	}

	var msg string
	switch lang {
	case "vn":
		msg = t.DoSayHelloInVn(name)
	case "en":
		msg = t.DoSayHelloInEn(name)
	default:
		http.Error(w, fmt.Sprintf("%s is not supported", lang), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"msg": msg})
}

func main() {
	r := chi.NewRouter()
	api := &TestApi{}

	r.Get("/hello/{name}/{langCode}", api.Hello)
	r.Post("/users", CreateUser)
	http.ListenAndServe(":8082", r)
}
