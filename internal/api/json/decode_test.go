package json_test

import (
	"live2text/internal/api/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

//nolint:gocognit // TODO: simplify?
func TestDecode(t *testing.T) {
	tests := []struct {
		name string

		contentType string
		body        string

		expected     testStruct
		expectError  bool
		expectedCode int
		expectMsg    string
	}{
		{
			name:         "successful decoding",
			contentType:  "application/json",
			body:         `{"name":"test","value":123}`,
			expected:     testStruct{"test", 123},
			expectError:  false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid content type",
			contentType:  "text/plain",
			body:         `{"name":"test","value":123}`,
			expected:     testStruct{},
			expectError:  true,
			expectedCode: http.StatusBadRequest,
			expectMsg:    "cannot decode content type",
		},
		{
			name:         "invalid JSON",
			contentType:  "application/json",
			body:         `{"name":"test","value":invalid}`,
			expected:     testStruct{},
			expectError:  true,
			expectedCode: http.StatusBadRequest,
			expectMsg:    "cannot decode request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			result, hasError := json.Decode[testStruct](w, req)

			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got %v", tt.expectError, hasError)
			}
			if result != tt.expected {
				t.Errorf("Expected result: %+v, got %+v", tt.expected, result)
			}
			if tt.expectError {
				if w.Code != tt.expectedCode {
					t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
				}
				if !strings.Contains(w.Body.String(), tt.expectMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.expectMsg, w.Body.String())
				}
			}
		})
	}
}
