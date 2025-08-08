//nolint:dupl
package btt

import (
	"net/http"

	"github.com/roboloop/live2text/internal/api/json"
	"github.com/roboloop/live2text/internal/api/validation"
	"github.com/roboloop/live2text/internal/services/btt"
)

type selectViewModeRequest struct {
	ViewMode string `json:"view_mode"`
}

func (r selectViewModeRequest) validate() map[string]string {
	problems := make(map[string]string)

	if r.ViewMode != string(btt.ViewModeClean) && r.ViewMode != string(btt.ViewModeEmbed) {
		problems["view_mode"] = "view_mode is not valid"
	}

	return problems
}

func (s *Server) SelectViewMode(w http.ResponseWriter, r *http.Request) {
	request, err := json.Decode[selectViewModeRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if problems := request.validate(); len(problems) > 0 {
		validation.Error(w, problems)
		return
	}

	err = s.services.Btt().SelectViewMode(r.Context(), btt.ViewMode(request.ViewMode))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
