package validation

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, problems map[string]string) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(problems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(buf.Bytes())
}
