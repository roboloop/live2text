package json_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	encodingjson "encoding/json"

	"live2text/internal/api/json"
)

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestDecode(t *testing.T) {
	t.Parallel()

	t.Run("unsupported content type", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		req.Header.Set("Content-Type", "text/plain")

		_, err := json.Decode[testStruct](req)

		require.ErrorContains(t, err, "cannot decode content type 'text/plain'")
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("invalid"))
		req.Header.Set("Content-Type", "application/json")

		_, err := json.Decode[testStruct](req)

		require.ErrorContains(t, err, "cannot decode request")
		var syntaxErr *encodingjson.SyntaxError
		require.ErrorAs(t, err, &syntaxErr)
	})

	t.Run("successful decoding", func(t *testing.T) {
		t.Parallel()

		body := `{"name":"foo","value":42}`
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		decoded, err := json.Decode[testStruct](req)

		require.NoError(t, err)
		require.Equal(t, testStruct{"foo", 42}, decoded)
	})
}
