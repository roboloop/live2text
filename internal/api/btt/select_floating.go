//nolint:dupl
package btt

import (
	"net/http"

	"live2text/internal/api/json"
	"live2text/internal/api/validation"
	"live2text/internal/services/btt"
)

type selectFloatingRequest struct {
	Floating string `json:"floating"`
}

func (r selectFloatingRequest) validate() map[string]string {
	problems := make(map[string]string)

	if r.Floating != string(btt.FloatingShown) && r.Floating != string(btt.FloatingHidden) {
		problems["floating"] = "floating is not valid"
	}

	return problems
}

func (s *Server) SelectFloating(w http.ResponseWriter, r *http.Request) {
	request, err := json.Decode[selectFloatingRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if problems := request.validate(); len(problems) > 0 {
		validation.Error(w, problems)
		return
	}

	err = s.services.Btt().SelectFloating(r.Context(), btt.Floating(request.Floating))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
