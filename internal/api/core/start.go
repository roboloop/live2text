package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/gordonklaus/portaudio"

	"live2text/internal/api/json"
	"live2text/internal/api/validation"
	"live2text/internal/services"
	"live2text/internal/services/recognition"
)

type startRequest struct {
	Device   string `json:"device"`
	Language string `json:"language"`
}

type startResponse struct {
	ID         string `json:"id"`
	SocketPath string `json:"socketPath"`
}

func (r startRequest) Valid(_ context.Context, s services.Services) (map[string]string, error) {
	problems := make(map[string]string)

	devices, err := s.Audio().List()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of devices: %w", err)
	}
	if !slices.ContainsFunc(devices, func(device *portaudio.DeviceInfo) bool {
		return device.Name == r.Device
	}) {
		problems["device"] = "device not found"
	}

	if !validation.IsValidLanguageCode(r.Language) {
		problems["language"] = "language is not valid"
	}

	return problems, nil
}

func (s *Server) Start(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request startRequest
	var responded bool
	if request, responded = json.Decode[startRequest](w, r); responded {
		return
	}
	if responded = validation.Validate(request, s.services, w, r); responded {
		return
	}

	id, socketPath, err := s.services.Recognition().Start(r.Context(), request.Device, request.Language)
	if err != nil {
		if errors.Is(err, recognition.ErrDeviceIsBusy) {
			json.Encode(errorResponse{err.Error()}, w, http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode(startResponse{id, socketPath}, w, http.StatusOK)
}
