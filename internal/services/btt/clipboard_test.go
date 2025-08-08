package btt_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/btt"
	"github.com/roboloop/live2text/internal/services/btt/client"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestSelectClipboard(t *testing.T) {
	t.Parallel()

	cc := setupClipboardComponent(
		t,
		func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock) {
			st.SelectSettingsMock.Expect(minimock.AnyContext, "Selected Clipboard", "LIVE2TEXT_SELECTED_CLIPBOARD", "foo").
				Return(nil)
		},
		nil,
	)
	err := cc.SelectClipboard(t.Context(), "foo")
	require.NoError(t, err)
}

func TestSelectedClipboard(t *testing.T) {
	t.Parallel()

	cc := setupClipboardComponent(
		t,
		func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock) {
			st.SelectedSettingMock.Expect(minimock.AnyContext, "LIVE2TEXT_SELECTED_CLIPBOARD").Return("foo", nil)
		},
		nil,
	)
	clipboard, err := cc.SelectedClipboard(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, "foo", clipboard)
}

func TestShowClipboard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock)
		expectErr  string
	}{
		{
			name: "cannot get selected clipboard",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock) {
				st.SelectedSettingMock.Return("", errors.New("something happened"))
			},
			expectErr: "cannot get selected clipboard",
		},
		{
			name: "cannot update clipboards",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock) {
				st.SelectedSettingMock.Return("Shown", nil)
				c.UpdateTriggerMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot update clipboards",
		},
		{
			name: "shown disabled",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock) {
				st.SelectedSettingMock.Return("", nil)
			},
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock) {
				st.SelectedSettingMock.Expect(minimock.AnyContext, "LIVE2TEXT_SELECTED_CLIPBOARD").Return("Shown", nil)
				c.UpdateTriggerMock.Times(2).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cc := setupClipboardComponent(t, tt.setupMocks, nil)

			err := cc.ShowClipboard(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHideClipboard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock)
		expectErr  string
	}{
		{
			name: "cannot update clipboard",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock) {
				c.UpdateTriggerMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot update clipboard",
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock) {
				c.UpdateTriggerMock.Times(2).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cc := setupClipboardComponent(t, tt.setupMocks, nil)

			err := cc.HideClipboard(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func setupClipboardComponent(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, c *client.ClientMock, st *btt.SettingsComponentMock),
	l *slog.Logger,
) btt.ClipboardComponent {
	mc := minimock.NewController(t)
	c := client.NewClientMock(mc)
	st := btt.NewSettingsComponentMock(mc)

	if setupMocks != nil {
		setupMocks(mc, c, st)
	}
	if l == nil {
		l = logger.NilLogger
	}

	return btt.NewClipboardComponent(l, c, st)
}
