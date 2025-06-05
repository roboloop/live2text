package btt

import (
	"net/http"

	"live2text/internal/api/json"
)

func (s *Server) ToggleListening(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := s.services.Btt().ToggleListening(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
