package btt

import (
	"live2text/internal/api/json"
	"net/http"
)

func (s *Server) LoadDevices(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := s.services.Btt().LoadDevices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
