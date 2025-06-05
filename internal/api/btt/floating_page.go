package btt

import (
	"net/http"
)

func (s *Server) FloatingPage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	page, err := s.services.Btt().Page()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte(page)); err != nil {
		s.logger.ErrorContext(r.Context(), "Failed to write response", "error", err)
	}
}
