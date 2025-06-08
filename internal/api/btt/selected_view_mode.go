package btt

import "net/http"

func (s *Server) SelectedViewMode(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	language, err := s.services.Btt().SelectedViewMode(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte(language)); err != nil {
		s.logger.ErrorContext(r.Context(), "Failed to write response", "error", err)
	}
}
