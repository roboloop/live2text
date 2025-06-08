package btt

import (
	"net/http"

	"live2text/internal/api/json"
)

func (s *Server) IsRunning(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	isRunning, err := s.services.Btt().IsRunning(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode(isRunning, w, http.StatusOK)
}
