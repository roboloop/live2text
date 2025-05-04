package encoding_test

import (
	"io"
	"live2text/internal/api/encoding"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	t.Run("Valid JSON payload", func(t *testing.T) {
		type testResponse struct {
			Foo string `json:"foo"`
			Bar int    `json:"bar"`
		}

		payload := &testResponse{
			Foo: "value",
			Bar: 42,
		}
		recorder := httptest.NewRecorder()
		err := encoding.Encode(recorder, http.StatusCreated, payload)
		if err != nil {
			t.Fatalf("Encode() returned error: got %v, expected nil", err)
		}
		httpResponse := recorder.Result()
		bodyBytes, _ := io.ReadAll(httpResponse.Body)
		expectedBody := "{\"foo\":\"value\",\"bar\":42}\n"

		if string(bodyBytes) != expectedBody {
			t.Errorf("Response body mismatch: got %v, expected %v", string(bodyBytes), expectedBody)
			return
		}

		if httpResponse.StatusCode != http.StatusCreated {
			t.Errorf("Status code mismatch: got %v, expected %v", httpResponse.StatusCode, http.StatusCreated)
			return
		}

		contentType := httpResponse.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type header mismatch: got %v, expected %v", contentType, "application/json")
			return
		}
	})

	t.Run("Cannot encode response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		err := encoding.Encode(recorder, http.StatusOK, func() {})
		if err == nil || !strings.Contains(err.Error(), "cannot encode response: ") {
			t.Fatalf("Encode() returned no error")
		}
	})
}

func TestDecode(t *testing.T) {
	t.Run("Valid JSON payload", func(t *testing.T) {
		type testRequest struct {
			Foo string `json:"foo"`
			Bar int    `json:"bar"`
		}
		bodyBytes := `{"foo":"value","bar":42}`
		httpRequest := httptest.NewRequest("GET", "/api", strings.NewReader(bodyBytes))
		httpRequest.Header.Set("Content-Type", "application/json")

		payload, err := encoding.Decode[testRequest](httpRequest)
		if err != nil {
			t.Fatalf("Decode() returned unxepected error: %v", err)
		}

		expectedPayload := testRequest{Foo: "value", Bar: 42}
		if *payload != expectedPayload {
			t.Fatalf("Decoded payload mismatch: got %v, expected %v", payload, expectedPayload)
		}
	})

	t.Run("Unsupported Content-Type", func(t *testing.T) {
		httpRequest := httptest.NewRequest("GET", "/api", strings.NewReader(""))
		httpRequest.Header.Set("Content-Type", "application/xml")

		_, err := encoding.Decode[string](httpRequest)
		expectedErr := "cannot decode content type 'application/xml'"
		if err == nil || err.Error() != expectedErr {
			t.Fatalf("Expected error: got %v, expected %v", err, expectedErr)
		}
	})

	t.Run("Invalid JSON payload", func(t *testing.T) {
		httpRequest := httptest.NewRequest("GET", "/api", strings.NewReader("foo"))
		httpRequest.Header.Set("Content-Type", "application/json")

		_, err := encoding.Decode[string](httpRequest)
		expectedErr := "cannot decode request"
		if err == nil || !strings.Contains(err.Error(), expectedErr) {
			t.Errorf("Expected error: got %v, expected contains %v", err, expectedErr)
			return
		}
	})
}
