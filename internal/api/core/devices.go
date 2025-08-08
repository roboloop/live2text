package core

import (
	"net/http"

	"github.com/roboloop/live2text/internal/api/json"
)

type devicesResponse struct {
	Devices []string `json:"devices"`
}

func (s *Server) Devices(w http.ResponseWriter, _ *http.Request) {
	devices, err := s.services.Audio().ListOfNames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode(&devicesResponse{devices}, w, http.StatusOK)
}
