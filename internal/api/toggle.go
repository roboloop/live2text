package api

import (
	"net/http"
)

//type toggleRequest struct {
//}

//type toggleResponse struct {
//	Id         string `json:"id"`
//	SocketPath string `json:"socketPath"`
//}

//func (r toggleRequest) Valid(_ context.Context, _ *Server) (map[string]string, error) {
//	return nil, nil
//}

func (s *Server) ToggleListening(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := s.services.Btt().ToggleListening(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encode("ok", w, http.StatusOK)
}
