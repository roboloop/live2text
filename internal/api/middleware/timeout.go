package middleware

import (
	"net/http"
	"time"
)

// TimeoutMiddleware timeouts any requests.
func TimeoutMiddleware(
	next http.Handler,
	timeout time.Duration,
	excludedPaths []string,
	isDebug bool,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isDebug {
			next.ServeHTTP(w, r)
			return
		}

		for _, p := range excludedPaths {
			if r.URL.Path == p {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.TimeoutHandler(next, timeout, "Request timed out").ServeHTTP(w, r)
	}
}
