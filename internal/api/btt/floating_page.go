package btt

import (
	"net/http"
)

func (s *Server) FloatingPage(w http.ResponseWriter, _ *http.Request) {
	page := s.services.Btt().FloatingPage()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(page))
}
