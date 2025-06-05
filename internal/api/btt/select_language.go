package btt

import (
	"context"
	"net/http"

	"live2text/internal/api/json"
	"live2text/internal/api/validation"
)

type selectLanguageRequest struct {
	Language string `json:"language"`
}

func (r selectLanguageRequest) Valid(context.Context, *Server) (map[string]string, error) {
	problems := make(map[string]string)

	if !validation.IsValidLanguageCode(r.Language) {
		problems["language"] = "language is not valid"
	}

	return problems, nil
}

func (s *Server) SelectLanguage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request, responded := json.Decode[selectLanguageRequest](w, r)
	if responded {
		return
	}

	err := s.services.Btt().SelectLanguage(r.Context(), request.Language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Encode("ok", w, http.StatusOK)
}
