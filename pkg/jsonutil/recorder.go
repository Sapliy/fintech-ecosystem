package jsonutil

import (
	"bytes"
	"net/http"
)

// ResponseRecorder captures the status code and body written to it.
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       bytes.Buffer
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}
