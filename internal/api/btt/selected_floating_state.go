package btt

import "net/http"

func (s *Server) SelectedFloatingState(w http.ResponseWriter, r *http.Request) {
	state, err := s.services.Btt().SelectedFloating(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(state))
}
