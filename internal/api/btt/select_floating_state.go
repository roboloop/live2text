package btt

import (
	"context"
	"net/http"

	"live2text/internal/api/json"
	"live2text/internal/services/btt"
)

type selectFloatingStateRequest struct {
	FloatingState string `json:"floating_state"`
}

func (r selectFloatingStateRequest) Valid(context.Context, *Server) (map[string]string, error) {
	problems := make(map[string]string)

	if r.FloatingState != btt.FloatingStateShown && r.FloatingState != btt.FloatingStateHidden {
		problems["floating_state"] = "floating_state is not valid"
	}

	return problems, nil
}

func (s *Server) SelectFloatingState(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request, responded := json.Decode[selectFloatingStateRequest](w, r)
	if responded {
		return
	}

	err := s.services.Btt().SelectFloatingState(r.Context(), btt.FloatingState(request.FloatingState))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
