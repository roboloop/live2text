package btt

import (
	"net/http"

	"github.com/roboloop/live2text/internal/api/json"
)

func (s *Server) IsRunning(w http.ResponseWriter, r *http.Request) {
	isRunning, err := s.services.Btt().IsRunning(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode(isRunning, w, http.StatusOK)
}
