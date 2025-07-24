package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	httpnet "net/http"

	"live2text/internal/services/btt/client/http"
	"live2text/internal/utils/logger"
)

type roundTripperFunc func(request *httpnet.Request) (*httpnet.Response, error)

func (rt roundTripperFunc) RoundTrip(r *httpnet.Request) (*httpnet.Response, error) {
	return rt(r)
}

type failingWriter struct{}

func (w *failingWriter) Read([]byte) (int, error) {
	return 0, errors.New("dummy error")
}

func (w *failingWriter) Close() error {
	return nil
}

func TestSend(t *testing.T) {
	t.Parallel()

	t.Run("cannot encode payload", func(t *testing.T) {
		t.Parallel()

		client := http.NewClient(logger.NilLogger, "", nil)

		body, err := client.Send(t.Context(), "", map[string]any{"foo": func() {}}, nil)

		require.Nil(t, body)
		require.ErrorContains(t, err, "cannot encode payload")
		var errSyntax *json.UnsupportedTypeError
		require.ErrorAs(t, err, &errSyntax)
	})

	t.Run("cannot create new request", func(t *testing.T) {
		t.Parallel()

		rt := func(req *httpnet.Request) (*httpnet.Response, error) {
			return &httpnet.Response{}, nil
		}
		httpClient := &httpnet.Client{Transport: roundTripperFunc(rt)}
		client := http.NewClient(logger.NilLogger, string([]byte{0x00}), httpClient)

		body, err := client.Send(t.Context(), "foo", nil, nil)

		require.Nil(t, body)
		require.ErrorContains(t, err, "cannot create new request")
	})

	t.Run("cannot send request", func(t *testing.T) {
		rt := func(req *httpnet.Request) (*httpnet.Response, error) {
			return nil, errors.New("dummy error")
		}
		httpClient := &httpnet.Client{Transport: roundTripperFunc(rt)}
		client := http.NewClient(logger.NilLogger, "foo", httpClient)

		body, err := client.Send(t.Context(), "foo", nil, nil)

		require.Nil(t, body)
		require.ErrorContains(t, err, "cannot send request")
		require.ErrorContains(t, err, "dummy error")
	})

	t.Run("unexpected response status code", func(t *testing.T) {
		t.Parallel()

		rt := func(req *httpnet.Request) (*httpnet.Response, error) {
			w := httptest.NewRecorder()
			w.WriteHeader(httpnet.StatusInternalServerError)
			return w.Result(), nil
		}
		httpClient := &httpnet.Client{Transport: roundTripperFunc(rt)}
		client := http.NewClient(logger.NilLogger, "foo", httpClient)

		body, err := client.Send(t.Context(), "", nil, nil)

		require.Nil(t, body)
		require.ErrorContains(t, err, "unexpected response status code 500")
	})

	t.Run("cannot read body", func(t *testing.T) {
		t.Parallel()

		rt := func(req *httpnet.Request) (*httpnet.Response, error) {
			resp := httptest.NewRecorder().Result()
			resp.Body = &failingWriter{}
			return resp, nil
		}
		httpClient := &httpnet.Client{Transport: roundTripperFunc(rt)}
		client := http.NewClient(logger.NilLogger, "foo", httpClient)

		body, err := client.Send(t.Context(), "", nil, nil)

		require.Nil(t, body)
		require.ErrorContains(t, err, "cannot read body")
		require.ErrorContains(t, err, "dummy error")
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		rt := func(req *httpnet.Request) (*httpnet.Response, error) {
			u, err := url.ParseQuery(req.URL.RawQuery)
			require.NoError(t, err)
			require.Equal(t, "bar", u.Get("foo"))
			require.JSONEq(t, "{\"key\":\"value\"}\n", u.Get("json"))
			require.Equal(t, httpnet.MethodGet, req.Method)
			require.Equal(t, "example.com", req.URL.Host)
			require.Equal(t, "/method_name/", req.URL.Path)

			w := httptest.NewRecorder()
			w.Body = bytes.NewBufferString("sample text")
			return w.Result(), nil
		}
		httpClient := &httpnet.Client{Transport: roundTripperFunc(rt)}
		client := http.NewClient(logger.NilLogger, "example.com", httpClient)

		body, err := client.Send(
			t.Context(),
			"method_name",
			map[string]any{"key": "value"},
			map[string]string{"foo": "bar"},
		)

		require.NoError(t, err)
		require.Equal(t, []byte("sample text"), body)
	})
}
