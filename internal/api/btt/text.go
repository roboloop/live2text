package btt

import (
	"net/http"
)

func (s *Server) Text(w http.ResponseWriter, r *http.Request) {
	text, err := s.services.Btt().Text(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(text))
}
