package json_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"live2text/internal/api/json"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name         string
		input        any
		status       int
		expectedCode int
		expectedJSON string
	}{
		{
			name: "successful encoding",
			input: struct {
				Name  string `json:"name"`
				Value int    `json:"value"`
			}{
				Name: "test", Value: 123,
			},
			status:       http.StatusOK,
			expectedCode: http.StatusOK,
			expectedJSON: `{"name":"test","value":123}`,
		},
		{
			name:         "nil pointer encoding",
			input:        (*struct{})(nil),
			status:       http.StatusOK,
			expectedCode: http.StatusOK,
			expectedJSON: `null`,
		},
		{
			name: "map encoding",
			input: map[string]any{
				"name":  "test",
				"value": 123,
			},
			status:       http.StatusCreated,
			expectedCode: http.StatusCreated,
			expectedJSON: `{"name":"test","value":123}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			json.Encode(tt.input, w, tt.status)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}
			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Expected content type 'application/json', got '%s'", contentType)
			}
			if tt.expectedJSON != strings.TrimSpace(w.Body.String()) {
				t.Errorf("JSON not equal:\nExpected: %v\nActual: %v", tt.expectedJSON, w.Body.String())
			}
		})
	}
}
