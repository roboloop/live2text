package btt

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/roboloop/live2text/internal/api/json"
	"github.com/roboloop/live2text/internal/api/validation"
	"github.com/roboloop/live2text/internal/services"
)

type selectDeviceRequest struct {
	Device string `json:"device"`
}

func (r selectDeviceRequest) validate(s services.Services) (map[string]string, error) {
	problems := make(map[string]string)

	devices, err := s.Audio().ListOfNames()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of devices: %w", err)
	}
	if !slices.Contains(devices, r.Device) {
		problems["device"] = "device not found"
	}

	return problems, nil
}

func (s *Server) SelectDevice(w http.ResponseWriter, r *http.Request) {
	request, err := json.Decode[selectDeviceRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	problems, err := request.validate(s.services)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(problems) > 0 {
		validation.Error(w, problems)
		return
	}

	err = s.services.Btt().SelectDevice(r.Context(), request.Device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
