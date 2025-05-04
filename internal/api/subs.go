package api

import (
	"context"
	"errors"
	"live2text/internal/services/recognition"
	"net/http"
)

type subsRequest struct {
	Id string `json:"id"`
}

func (r subsRequest) Valid(_ context.Context, api *Server) (map[string]string, error) {
	return nil, nil
}

func (s *Server) Subs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	request, responded := decode[subsRequest](s, w, r)
	if responded {
		return
	}

	text, err := s.services.Recognition().Subs(r.Context(), request.Id)
	if err != nil {
		if errors.Is(err, recognition.NoTaskError) {
			encode(errorResponse{err.Error()}, w, http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Avoid a json response, write just a simple text
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}
