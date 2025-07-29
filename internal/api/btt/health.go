package btt

import (
	"net/http"
)

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	ok := s.services.Btt().Health(r.Context())
	if !ok {
		http.Error(w, "btt is not running", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
