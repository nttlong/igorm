package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// Payload test
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

// go test -bench ^BenchmarkChiPostUser$ -benchmem -cpuprofile cpu_chi.out -memprofile mem_chi.out
func BenchmarkChiPostUser(b *testing.B) {
	router := chi.NewRouter()
	router.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		var u User
		_ = json.NewDecoder(r.Body).Decode(&u)
		w.WriteHeader(http.StatusOK)
	})

	body, _ := json.Marshal(UserInput{
		Name: "abc",
		Age:  30,
		Address: Address{
			City: "Hanoi",
		},
	})
	req := httptest.NewRequest("POST", "/user", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}
