package btt

import "net/http"

func (s *Server) SelectedFloatingState(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	state, err := s.services.Btt().SelectedFloatingState(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte(state)); err != nil {
		s.logger.ErrorContext(r.Context(), "Failed to write response", "error", err)
	}
}
