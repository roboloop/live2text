package validation

import (
	"bytes"
	"context"
	"encoding/json"
	"live2text/internal/services"
	"net/http"
)

type Validator interface {
	Valid(context.Context, services.Services) (map[string]string, error)
}

func Validate(v Validator, s services.Services, w http.ResponseWriter, r *http.Request) bool {
	problems, err := v.Valid(r.Context(), s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	if len(problems) > 0 {
		var buf = bytes.Buffer{}
		if err = json.NewEncoder(&buf).Encode(problems); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return true
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write(buf.Bytes())
		return true
	}

	return false
}
