package btt_test

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/btt"
	"github.com/roboloop/live2text/internal/services/btt/client"
)

func TestSelectViewMode(t *testing.T) {
	t.Parallel()

	vmc := setupViewModeComponent(
		t,
		func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock) {
			s.SelectSettingsMock.Expect(minimock.AnyContext, "Selected View Mode", "LIVE2TEXT_SELECTED_VIEW_MODE", "foo").
				Return(nil)
		},
	)
	err := vmc.SelectViewMode(t.Context(), "foo")
	require.NoError(t, err)
}

func TestSelectedViewMode(t *testing.T) {
	t.Parallel()

	vmc := setupViewModeComponent(
		t,
		func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock) {
			s.SelectedSettingMock.Expect(minimock.AnyContext, "LIVE2TEXT_SELECTED_VIEW_MODE").Return("foo", nil)
		},
	)
	viewMode, err := vmc.SelectedViewMode(t.Context())
	require.NoError(t, err)
	require.EqualValues(t, "foo", viewMode)
}

func TestEnableCleanMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock)
		expectErr  string
	}{
		{
			name: "cannot get selected view mode",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock) {
				s.SelectedSettingMock.Return("", errors.New("something happened"))
			},
			expectErr: "cannot get selected view mode",
		},
		{
			name: "cannot open dir",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock) {
				s.SelectedSettingMock.Return("Clean", nil)
				c.TriggerActionMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot open dir",
		},
		{
			name: "embed view mode",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock) {
				s.SelectedSettingMock.Return("Embed", nil)
			},
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock) {
				s.SelectedSettingMock.Expect(minimock.AnyContext, "LIVE2TEXT_SELECTED_VIEW_MODE").Return("Clean", nil)
				c.TriggerActionMock.Return(nil)
			},
			expectErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			vmc := setupViewModeComponent(t, tt.setupMocks)

			err := vmc.EnableCleanMode(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDisableCleanView(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock)
		expectErr  string
	}{
		{
			name: "cannot close dir",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock) {
				c.TriggerActionMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot close dir",
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock) {
				c.TriggerActionMock.Return(nil)
			},
			expectErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			vmc := setupViewModeComponent(t, tt.setupMocks)

			err := vmc.DisableCleanView(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func setupViewModeComponent(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, c *client.ClientMock, s *btt.SettingsComponentMock),
) btt.ViewModeComponent {
	mc := minimock.NewController(t)
	c := client.NewClientMock(mc)
	s := btt.NewSettingsComponentMock(mc)

	if setupMocks != nil {
		setupMocks(mc, c, s)
	}

	return btt.NewViewModeComponent(c, s)
}
