package middleware

import (
	"fmt"
	"net/http"
)

//type ErrorMiddleware struct {
//	next http.Handler
//}
//
//func NewErrorMiddleware(next http.Handler) http.Handler {
//	return &ErrorMiddleware{next}
//}
//
//func (m *ErrorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	defer func() {
//		if err := recover(); err != nil {
//			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//		}
//	}()
//
//	m.next.ServeHTTP(w, r)
//}

func ErrorMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}
}
