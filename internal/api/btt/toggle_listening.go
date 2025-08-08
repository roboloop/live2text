package btt

import (
	"errors"
	"net/http"

	"github.com/roboloop/live2text/internal/api/json"
	"github.com/roboloop/live2text/internal/services/btt"
)

func (s *Server) ToggleListening(w http.ResponseWriter, r *http.Request) {
	if err := s.services.Btt().ToggleListening(r.Context()); err != nil {
		if errors.Is(err, btt.ErrDeviceNotSelected) ||
			errors.Is(err, btt.ErrDeviceIsUnavailable) ||
			errors.Is(err, btt.ErrLanguageNotSelected) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
