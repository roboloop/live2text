package core

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/roboloop/live2text/internal/api/json"
	"github.com/roboloop/live2text/internal/api/validation"
	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/recognition"
)

type startRequest struct {
	Device   string `json:"device"`
	Language string `json:"language"`
}

type startResponse struct {
	ID         string `json:"id"`
	SocketPath string `json:"socketPath"`
}

func (r startRequest) validate(s services.Services) (map[string]string, error) {
	problems := make(map[string]string)

	devices, err := s.Audio().ListOfNames()
	if err != nil {
		return nil, fmt.Errorf("cannot get list of devices: %w", err)
	}
	if !slices.Contains(devices, r.Device) {
		problems["device"] = "device not found"
	}

	if !validation.IsValidLanguageCode(r.Language) {
		problems["language"] = "language is not valid"
	}

	return problems, nil
}

func (s *Server) Start(w http.ResponseWriter, r *http.Request) {
	request, err := json.Decode[startRequest](r)
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
