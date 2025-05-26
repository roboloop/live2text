package core

import (
	"errors"
	"live2text/internal/api/json"
	"live2text/internal/services/recognition"
	"net/http"
)

type stopRequest struct {
	ID string `json:"id"`
}

func (s *Server) Stop(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request, responded := json.Decode[stopRequest](w, r)
	if responded {
		return
	}

	err := s.services.Recognition().Stop(r.Context(), request.ID)
	if err != nil {
		if errors.Is(err, recognition.ErrNoDeviceBusy) {
			json.Encode(errorResponse{err.Error()}, w, http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
