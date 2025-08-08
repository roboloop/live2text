package btt

import (
	"net/http"

	"github.com/roboloop/live2text/internal/api/json"
)

func (s *Server) LoadDevices(w http.ResponseWriter, r *http.Request) {
	err := s.services.Btt().LoadDevices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
