package core

import (
	"net/http"

	"live2text/internal/api/json"
)

type healthResponse string

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	json.Encode(healthResponse("ok"), w, http.StatusOK)
}
