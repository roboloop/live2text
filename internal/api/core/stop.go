package core

import (
	"errors"
	"net/http"

	"github.com/roboloop/live2text/internal/api/json"
	"github.com/roboloop/live2text/internal/services/recognition"
)

type stopRequest struct {
	ID string `json:"id"`
}

func (s *Server) Stop(w http.ResponseWriter, r *http.Request) {
	request, err := json.Decode[stopRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.services.Recognition().Stop(r.Context(), request.ID)
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
