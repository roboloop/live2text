package core

import (
	"errors"
	"live2text/internal/api/json"
	"live2text/internal/services/recognition"
	"net/http"
)

type subsRequest struct {
	Id string `json:"id"`
}

func (s *Server) Subs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	request, responded := json.Decode[subsRequest](w, r)
	if responded {
		return
	}

	text, err := s.services.Recognition().Subs(r.Context(), request.Id)
	if err != nil {
		if errors.Is(err, recognition.NoTaskError) {
			json.Encode(errorResponse{err.Error()}, w, http.StatusBadRequest)
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
