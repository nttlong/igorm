package handlers

import (
	"bytes"
	"net/http"
)

type MockResponseWriter struct {
	HeaderMap http.Header
	Body      bytes.Buffer
	Status    int
}

func newMockResponseWriter() http.ResponseWriter {
	return &MockResponseWriter{
		HeaderMap: make(http.Header),
	}
}

func (m *MockResponseWriter) Header() http.Header {
	return m.HeaderMap
}

func (m *MockResponseWriter) Write(b []byte) (int, error) {
	return m.Body.Write(b)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Status = statusCode
}
