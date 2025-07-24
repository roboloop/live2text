package btt

import (
	"net/http"

	"live2text/internal/api/json"
	"live2text/internal/api/validation"
)

type selectLanguageRequest struct {
	Language string `json:"language"`
}

func (r selectLanguageRequest) validate() map[string]string {
	problems := make(map[string]string)

	if !validation.IsValidLanguageCode(r.Language) {
		problems["language"] = "language is not valid"
	}

	return problems
}

func (s *Server) SelectLanguage(w http.ResponseWriter, r *http.Request) {
	request, err := json.Decode[selectLanguageRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if problems := request.validate(); len(problems) > 0 {
		validation.Error(w, problems)
		return
	}

	err = s.services.Btt().SelectLanguage(r.Context(), request.Language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
