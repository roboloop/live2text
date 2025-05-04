package encoding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("cannot encode response: %w", err)
	}

	return nil
}

func Decode[T any](r *http.Request) (*T, error) {
	value := new(T)
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("cannot decode content type '%s'", contentType)
	}

	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		return nil, fmt.Errorf("cannot decode request: %w", err)
	}
	return value, nil
}
