package api

import "net/http"

type healthResponse string

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	encode(healthResponse("ok"), w, http.StatusOK)
}
