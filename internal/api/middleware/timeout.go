package middleware

import (
	"net/http"
	"time"
)

func TimeoutMiddleware(next http.Handler, timeout time.Duration, excludedPaths []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, p := range excludedPaths {
			if r.URL.Path == p {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.TimeoutHandler(next, timeout, "Request timed out").ServeHTTP(w, r)
	}
}
