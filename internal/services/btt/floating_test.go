package btt_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/btt"
	"github.com/roboloop/live2text/internal/services/btt/client"
	"github.com/roboloop/live2text/internal/services/btt/storage"
	"github.com/roboloop/live2text/internal/services/btt/tmpl"
	"github.com/roboloop/live2text/internal/services/recognition"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestSelectFloating(t *testing.T) {
	t.Parallel()

	fc := setupFloatingComponent(
		t,
		func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
			st.SelectSettingsMock.Expect(minimock.AnyContext, "Selected Floating", "LIVE2TEXT_SELECTED_FLOATING", "foo").
				Return(nil)
		},
		nil,
	)
	err := fc.SelectFloating(t.Context(), "foo")
	require.NoError(t, err)
}

func TestSelectedFloating(t *testing.T) {
	t.Parallel()

	fc := setupFloatingComponent(
		t,
		func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
			st.SelectedSettingMock.Expect(minimock.AnyContext, "LIVE2TEXT_SELECTED_FLOATING").Return("foo", nil)
		},
		nil,
	)
	floating, err := fc.SelectedFloating(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, "foo", floating)
}

func TestShowFloating(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock)
		expectErr  string
	}{
		{
			name: "cannot get selected floating",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
				st.SelectedSettingMock.Return("", errors.New("something happened"))
			},
			expectErr: "cannot get selected floating",
		},
		{
			name: "cannot show floating",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
				st.SelectedSettingMock.Return("Shown", nil)
				c.TriggerActionMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot show floating",
		},
		{
			name: "shown disabled",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
				st.SelectedSettingMock.Return("", nil)
			},
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
				st.SelectedSettingMock.Expect(minimock.AnyContext, "LIVE2TEXT_SELECTED_FLOATING").Return("Shown", nil)
				c.TriggerActionMock.Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fc := setupFloatingComponent(t, tt.setupMocks, nil)

			err := fc.ShowFloating(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHideFloating(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock)
		expectErr  string
	}{
		{
			name: "cannot hide floating",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
				c.TriggerActionMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot hide floating",
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
				c.TriggerActionMock.Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fc := setupFloatingComponent(t, tt.setupMocks, nil)

			err := fc.HideFloating(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestFloatingPage(t *testing.T) {
	t.Parallel()

	fc := setupFloatingComponent(
		t,
		func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
			r.FloatingPageMock.Return("foo")
		},
		nil,
	)

	page := fc.FloatingPage()
	require.Equal(t, "foo", page)
}

func TestStreamText(t *testing.T) {
	t.Parallel()

	t.Run("cannot get listening to socket variable", func(t *testing.T) {
		t.Parallel()

		fc := setupFloatingComponent(
			t,
			func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
				s.GetValueMock.Return("", errors.New("something happened"))
			},
			nil,
		)
		_, _, err := fc.StreamText(t.Context())
		require.Error(t, err)
		require.ErrorContains(t, err, "cannot get listening to socket variable")
	})

	t.Run("cannot get text", func(t *testing.T) {
		t.Parallel()

		fc := setupFloatingComponent(
			t,
			func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
				s.GetValueMock.Return("foo", nil)
				rg.TextMock.Return("", errors.New("something happened"))
			},
			nil,
		)

		_, errCh, err := fc.StreamText(t.Context())

		require.NoError(t, err)
		errText := <-errCh
		require.Error(t, errText)
		require.ErrorContains(t, errText, "something happened")
	})

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		fc := setupFloatingComponent(
			t,
			func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock) {
				s.GetValueMock.Return("foo", nil)
				rg.TextMock.Return("bar", nil)
			},
			nil,
		)

		ctx, cancel := context.WithCancel(t.Context())
		textCh, _, err := fc.StreamText(ctx)
		require.NoError(t, err)
		text := <-textCh
		require.Equal(t, "bar", text)
		cancel()
	})
}

func setupFloatingComponent(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, st *btt.SettingsComponentMock),
	l *slog.Logger,
) btt.FloatingComponent {
	mc := minimock.NewController(t)
	rg := recognition.NewRecognitionMock(mc)
	c := client.NewClientMock(mc)
	s := storage.NewStorageMock(mc)
	r := tmpl.NewRendererMock(mc)
	st := btt.NewSettingsComponentMock(mc)

	if setupMocks != nil {
		setupMocks(mc, rg, c, s, r, st)
	}
	if l == nil {
		l = logger.NilLogger
	}

	return btt.NewFloatingComponent(l, rg, c, s, r, st)
}
