package core

import (
	"live2text/internal/api/json"
	"net/http"
)

type healthResponse string

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	json.Encode(healthResponse("ok"), w, http.StatusOK)
}
