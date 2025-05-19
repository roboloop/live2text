package api

import "net/http"

func (s *Server) SelectedLanguage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	language, err := s.services.Btt().SelectedLanguage(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(language))
}
