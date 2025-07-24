//nolint:dupl
package btt

import (
	"net/http"

	"live2text/internal/api/json"
	"live2text/internal/api/validation"
	"live2text/internal/services/btt"
)

type selectFloatingStateRequest struct {
	FloatingState string `json:"floating_state"`
}

func (r selectFloatingStateRequest) validate() map[string]string {
	problems := make(map[string]string)

	if r.FloatingState != btt.FloatingShown && r.FloatingState != btt.FloatingHidden {
		problems["floating_state"] = "floating_state is not valid"
	}

	return problems
}

func (s *Server) SelectFloatingState(w http.ResponseWriter, r *http.Request) {
	request, err := json.Decode[selectFloatingStateRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if problems := request.validate(); len(problems) > 0 {
		validation.Error(w, problems)
		return
	}

	err = s.services.Btt().SelectFloating(r.Context(), btt.Floating(request.FloatingState))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
