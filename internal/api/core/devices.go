package core

import (
	"net/http"

	"live2text/internal/api/json"
)

type devicesResponse struct {
	Devices []string `json:"devices"`
}

func (s *Server) Devices(w http.ResponseWriter, _ *http.Request) {
	devices, err := s.services.Audio().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var deviceNames []string
	for _, device := range devices {
		deviceNames = append(deviceNames, device.Name)
	}

	json.Encode(&devicesResponse{deviceNames}, w, http.StatusOK)
}
