package json_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"live2text/internal/api/json"
)

func TestEncode(t *testing.T) {
	t.Parallel()

	t.Run("encode error", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		json.Encode(func() {}, w, http.StatusOK)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "unsupported type")
	})

	t.Run("successful encoding", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		json.Encode("ok", w, http.StatusOK)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), `"ok"`)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})
}
