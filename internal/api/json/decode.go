package json

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Decode[T any](r *http.Request) (T, error) {
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
