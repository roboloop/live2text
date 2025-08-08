package core

import (
	"errors"
	"net/http"

	"github.com/roboloop/live2text/internal/api/json"
	"github.com/roboloop/live2text/internal/services/recognition"
)

type textRequest struct {
	ID string `json:"id"`
}

func (s *Server) Text(w http.ResponseWriter, r *http.Request) {
	request, err := json.Decode[textRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	text, err := s.services.Recognition().Text(r.Context(), request.ID)
	if err != nil {
		if errors.Is(err, recognition.ErrNoTaskFound) {
			json.Encode(errorResponse{err.Error()}, w, http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Avoid a json response, write just a simple text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(text))
}
