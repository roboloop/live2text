package api

import (
	"bytes"
	"context"
	"encoding/json"
	"live2text/internal/api/encoding"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

type Validator interface {
	Valid(context.Context, *Server) (map[string]string, error)
}

func decode[T Validator](s *Server, w http.ResponseWriter, r *http.Request) (*T, bool) {
	req, err := encoding.Decode[T](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, true
	}

	problems, err := (*req).Valid(r.Context(), s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, true
	}
	if len(problems) > 0 {
		var buf = bytes.Buffer{}
		if err = json.NewEncoder(&buf).Encode(problems); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, true
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(buf.Bytes())
		//w.Write(buf.Re)
		return nil, true
	}

	return req, false
}

func encode[T any](v T, w http.ResponseWriter, status int) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(buf.Bytes()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
