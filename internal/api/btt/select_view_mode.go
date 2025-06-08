package btt

import (
	"context"
	"net/http"

	"live2text/internal/api/json"
	"live2text/internal/services/btt"
)

type selectViewModeRequest struct {
	ViewMode string `json:"view_mode"`
}

func (r selectViewModeRequest) Valid(context.Context, *Server) (map[string]string, error) {
	problems := make(map[string]string)

	if r.ViewMode != btt.ViewModeClean && r.ViewMode != btt.ViewModeEmbed {
		problems["view_mode"] = "view_mode is not valid"
	}

	return problems, nil
}

func (s *Server) SelectViewMode(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request, responded := json.Decode[selectViewModeRequest](w, r)
	if responded {
		return
	}

	err := s.services.Btt().SelectViewMode(r.Context(), btt.ViewMode(request.ViewMode))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
