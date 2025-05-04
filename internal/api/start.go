package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"live2text/internal/api/validation"
	"live2text/internal/services/recognition"
	"net/http"
	"slices"
)

type startRequest struct {
	Device   string `json:"device"`
	Language string `json:"language"`
}

type startResponse struct {
	Id         string `json:"id"`
	SocketPath string `json:"socketPath"`
}

func (r startRequest) Valid(_ context.Context, api *Server) (map[string]string, error) {
	problems := make(map[string]string)

	devices, err := api.services.Audio().List()
	if err != nil {
		return nil, fmt.Errorf("could not get list of devices: %w", err)
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
	request, responded := decode[startRequest](s, w, r)
	if responded {
		return
	}

	id, socketPath, err := s.services.Recognition().Start(r.Context(), request.Device, request.Language)
	if err != nil {
		if errors.Is(err, recognition.DeviceIsBusyError) {
			encode(errorResponse{err.Error()}, w, http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encode(&startResponse{id, socketPath}, w, http.StatusOK)
}
