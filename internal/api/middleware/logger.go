package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"time"
)

//type LoggerMiddleware struct {
//	next   http.Handler
//	logger *slog.Logger
//}
//
//func NewLoggerMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
//	return &LoggerMiddleware{next, logger}
//}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

//func (m *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	since := time.Now()
//	ww := &responseWriter{
//		ResponseWriter: w,
//	}
//
//	m.next.ServeHTTP(ww, r)
//
//	duration := time.Now().Sub(since)
//
//	level := slog.LevelInfo
//	if ww.status >= 500 {
//		level = slog.LevelError
//	} else if ww.status >= 400 {
//		level = slog.LevelInfo
//	}
//
//	host, _, err := net.SplitHostPort(r.RemoteAddr)
//	if err != nil {
//		host = r.RemoteAddr
//	}
//
//	m.logger.LogAttrs(r.Context(), level, "HTTP request",
//		slog.String("method", r.Method),
//		slog.String("path", r.URL.Path),
//		slog.Int("status", ww.status),
//		slog.String("duration", duration.String()),
//		slog.String("IP", host),
//		slog.String("user_agent", r.UserAgent()),
//	)
//}

func LoggerMiddleware(next http.Handler, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		since := time.Now()
		ww := &responseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(ww, r)

		duration := time.Now().Sub(since)

		level := slog.LevelInfo
		if ww.status >= 500 {
			level = slog.LevelError
		} else if ww.status >= 400 {
			level = slog.LevelInfo
		}

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		logger.LogAttrs(r.Context(), level, "HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", ww.status),
			slog.String("duration", duration.String()),
			slog.String("IP", host),
			slog.String("user_agent", r.UserAgent()),
		)
	}
}
