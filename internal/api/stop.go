package api

import (
	"context"
	"errors"
	"live2text/internal/services/recognition"
	"net/http"
)

type stopRequest struct {
	Id string `json:"id"`
}

func (r stopRequest) Valid(_ context.Context, api *Server) (map[string]string, error) {
	problems := make(map[string]string)

	//devices, err := api.services.Audio().List()
	//if err != nil {
	//	return nil, fmt.Errorf("cannot get list of devices: %w", err)
	//}
	//if !slices.ContainsFunc(devices, func(device *portaudio.DeviceInfo) bool {
	//	return device.Name == r.Device
	//}) {
	//	problems["device"] = "device not found"
	//}

	return problems, nil
}

func (s *Server) Stop(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request, responded := decode[stopRequest](s, w, r)
	if responded {
		return
	}

	err := s.services.Recognition().Stop(r.Context(), request.Id)
	if err != nil {
		if errors.Is(err, recognition.NoDeviceBusyError) {
			encode(errorResponse{err.Error()}, w, http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
