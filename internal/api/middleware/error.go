package middleware

import (
	"fmt"
	"net/http"
)

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
