package validation_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/api/validation"
)

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("responses with error", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()

		validation.Error(w, map[string]string{"foo": "bar"})

		resp := w.Result()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		require.Contains(t, w.Body.String(), `{"foo":"bar"}`)
	})
}
