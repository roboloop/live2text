package middleware

import (
	"log/slog"
	"net/http"
)

// BodyCloserMiddleware closes body request; according to the source code, it's totally unnecessary.
func BodyCloserMiddleware(next http.Handler, l *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				l.Error("Failed close body", "error", err)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
