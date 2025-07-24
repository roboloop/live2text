package btt_test

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"live2text/internal/services"
	"live2text/internal/services/btt"
	"live2text/internal/utils/logger"
)

type unsupportedFlusherWriter struct {
	w *httptest.ResponseRecorder
}

func (w *unsupportedFlusherWriter) Header() http.Header {
	return w.w.Header()
}

func (w *unsupportedFlusherWriter) Write(n []byte) (int, error) {
	return w.w.Write(n)
}

func (w *unsupportedFlusherWriter) WriteHeader(code int) {
	w.w.WriteHeader(code)
}

func TestTextStream(t *testing.T) {
	t.Parallel()

	t.Run("streaming unsupported", func(t *testing.T) {
		t.Parallel()

		server := setupServer(t, func(_ *minimock.Controller, s *services.ServicesMock) {}, nil)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := &unsupportedFlusherWriter{httptest.NewRecorder()}

		server.TextStream(w, req)

		require.Equal(t, http.StatusOK, w.w.Code)
		require.Contains(t, "event: failed\ndata: Streaming unsupported\n\n", w.w.Body.String())
	})

	t.Run("stream text failed", func(t *testing.T) {
		t.Parallel()

		server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
			b := btt.NewBttMock(mc)
			b.StreamTextMock.Return(nil, nil, errors.New("dummy error"))
			s.BttMock.Return(b)
		}, nil)

		w := performRequest(t, server.TextStream, "")

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, "event: failed\ndata: dummy error\n\n", w.Body.String())
	})

	t.Run("request cancelled", func(t *testing.T) {
		t.Parallel()

		textCh, errCh := make(chan string, 10), make(chan error, 1)
		defer close(textCh)
		defer close(errCh)
		l, h := logger.NewCaptureLogger()
		server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
			b := btt.NewBttMock(mc)
			b.StreamTextMock.Return(textCh, errCh, nil)
			s.BttMock.Return(b)
		}, l)

		ctx, cancel := context.WithCancel(t.Context())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			server.TextStream(w, req)
		}()
		cancel()
		wg.Wait()

		require.Equal(t, http.StatusOK, w.Code)
		require.Empty(t, w.Body.String())
		require.Len(t, h.Logs, 1)
		require.Equal(t, slog.LevelError, h.Logs[0].Level)
		require.Equal(t, "Request context cancelled", h.Logs[0].Msg)
	})

	t.Run("error during stream", func(t *testing.T) {
		t.Parallel()

		textCh, errCh := make(chan string, 10), make(chan error, 1)
		defer close(textCh)
		defer close(errCh)

		l, h := logger.NewCaptureLogger()
		server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
			b := btt.NewBttMock(mc)
			b.StreamTextMock.Return(textCh, errCh, nil)
			s.BttMock.Return(b)
		}, l)

		errCh <- errors.New("dummy error")
		w := performRequest(t, server.TextStream, "")

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, "event: failed\ndata: dummy error\n\n", w.Body.String())
		require.Len(t, h.Logs, 1)
		require.Equal(t, slog.LevelError, h.Logs[0].Level)
		require.Equal(t, "Error during stream", h.Logs[0].Msg)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		textCh, errCh := make(chan string, 10), make(chan error, 1)
		defer close(errCh)

		l, h := logger.NewCaptureLogger()
		server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
			b := btt.NewBttMock(mc)
			b.StreamTextMock.Return(textCh, errCh, nil)
			s.BttMock.Return(b)
		}, l)

		textCh <- "sample text"
		close(textCh)

		w := performRequest(t, server.TextStream, "")

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, "event: message\ndata: sample text\n\n", w.Body.String())
		require.Len(t, h.Logs, 1)
		require.Equal(t, slog.LevelInfo, h.Logs[0].Level)
		require.Equal(t, "TextStream channel closed", h.Logs[0].Msg)
	})
}
