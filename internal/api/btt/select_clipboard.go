//nolint:dupl
package btt

import (
	"net/http"

	"live2text/internal/api/json"
	"live2text/internal/api/validation"
	"live2text/internal/services/btt"
)

type selectClipboardRequest struct {
	Clipboard string `json:"clipboard"`
}

func (r selectClipboardRequest) validate() map[string]string {
	problems := make(map[string]string)

	if r.Clipboard != string(btt.ClipboardShown) && r.Clipboard != string(btt.ClipboardHidden) {
		problems["clipboard"] = "clipboard is not valid"
	}

	return problems
}

func (s *Server) SelectClipboard(w http.ResponseWriter, r *http.Request) {
	request, err := json.Decode[selectClipboardRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if problems := request.validate(); len(problems) > 0 {
		validation.Error(w, problems)
		return
	}

	err = s.services.Btt().SelectClipboard(r.Context(), btt.Clipboard(request.Clipboard))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
