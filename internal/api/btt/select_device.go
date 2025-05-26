package btt

import (
	"context"
	"fmt"
	"live2text/internal/api/json"
	"net/http"
	"slices"

	"github.com/gordonklaus/portaudio"
)

type selectDeviceRequest struct {
	Device string `json:"device"`
}

func (r selectDeviceRequest) Valid(_ context.Context, api *Server) (map[string]string, error) {
	problems := make(map[string]string)

	devices, err := api.services.Audio().List()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of devices: %w", err)
	}
	if !slices.ContainsFunc(devices, func(device *portaudio.DeviceInfo) bool {
		return device.Name == r.Device
	}) {
		problems["device"] = "device not found"
	}

	return problems, nil
}

func (s *Server) SelectDevice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request, responded := json.Decode[selectDeviceRequest](w, r)
	if responded {
		return
	}

	err := s.services.Btt().SelectDevice(r.Context(), request.Device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
