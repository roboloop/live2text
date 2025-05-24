package middleware

import (
	"bytes"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	error  bytes.Buffer
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(body []byte) (int, error) {
	if w.status >= http.StatusInternalServerError {
		w.error.Write(body)
	}

	return w.ResponseWriter.Write(body)
}

func LoggerMiddleware(next http.Handler, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		since := time.Now()
		ww := &responseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(ww, r)

		duration := time.Now().Sub(since)

		level := slog.LevelInfo
		if ww.status >= http.StatusInternalServerError {
			level = slog.LevelError
		} else if ww.status >= http.StatusBadRequest {
			level = slog.LevelInfo
		}

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", ww.status),
			slog.String("duration", duration.String()),
			slog.String("IP", host),
			slog.String("user_agent", r.UserAgent()),
		}

		if ww.status >= http.StatusInternalServerError {
			attrs = append(attrs, slog.String("error", ww.error.String()))
		}

		logger.LogAttrs(r.Context(), level, "HTTP request", attrs...)
	}
}
