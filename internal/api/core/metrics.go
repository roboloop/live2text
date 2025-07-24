package core

import (
	"net/http"
)

func (s *Server) Metrics(w http.ResponseWriter, _ *http.Request) {
	s.services.Metrics().WritePrometheus(w)
}
