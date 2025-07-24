package btt

import "net/http"

func (s *Server) SelectedLanguage(w http.ResponseWriter, r *http.Request) {
	language, err := s.services.Btt().SelectedLanguage(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(language))
}
