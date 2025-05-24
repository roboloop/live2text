package json

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Encode[T any](v T, w http.ResponseWriter, status int) {
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
