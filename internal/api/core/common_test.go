package core_test

import (
	"errors"
	"net/http"
)

type errorResponseWriter struct{}

func (e *errorResponseWriter) Header() http.Header {
	return http.Header{}
}

func (e *errorResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func (e *errorResponseWriter) WriteHeader(int) {
}
