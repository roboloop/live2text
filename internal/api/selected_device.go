package api

import "net/http"

func (s *Server) SelectedDevice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	device, err := s.services.Btt().SelectedDevice(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(device))
}
