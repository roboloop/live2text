package api

import (
	"net/http"
)

func (s *Server) Metrics(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s.services.Metrics().WritePrometheus(w)
}
