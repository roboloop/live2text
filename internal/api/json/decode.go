package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Decode[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	decoded, err := decode[T](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return decoded, true
	}

	return decoded, false
}

func decode[T any](r *http.Request) (T, error) {
	var value T
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return value, fmt.Errorf("cannot decode content type '%s'", contentType)
	}

	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		return value, fmt.Errorf("cannot decode request: %w", err)
	}
	return value, nil
}
