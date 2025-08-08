package btt_test

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/btt"
	"github.com/roboloop/live2text/internal/services/btt/client"
	"github.com/roboloop/live2text/internal/services/btt/client/trigger"
	"github.com/roboloop/live2text/internal/services/btt/storage"
)

func TestSelectSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *client.ClientMock, s *storage.StorageMock)
		title      trigger.Title
		key        storage.Key
		value      string
		expectErr  string
	}{
		{
			name: "cannot set value",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *storage.StorageMock) {
				s.SetValueMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot set value",
		},
		{
			name: "cannot refresh selected setting",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *storage.StorageMock) {
				s.SetValueMock.Return(nil)
				c.RefreshTriggerMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot refresh selected setting",
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *storage.StorageMock) {
				s.SetValueMock.Expect(minimock.AnyContext, "bar", "baz").Return(nil)
				c.RefreshTriggerMock.Expect(minimock.AnyContext, "foo").Return(nil)
			},
			title: "foo",
			key:   "bar",
			value: "baz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc := setupSettingsComponent(t, tt.setupMocks)

			err := sc.SelectSettings(t.Context(), tt.title, tt.key, tt.value)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestSelectedSetting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMocks  func(mc *minimock.Controller, c *client.ClientMock, s *storage.StorageMock)
		key         storage.Key
		expectValue string
		expectErr   string
	}{
		{
			name: "cannot get selected setting",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *storage.StorageMock) {
				s.GetValueMock.Return("", errors.New("something happened"))
			},
			expectErr: "cannot get selected setting",
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, s *storage.StorageMock) {
				s.GetValueMock.Expect(minimock.AnyContext, "foo").Return("bar", nil)
			},
			key:         "foo",
			expectValue: "bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sc := setupSettingsComponent(t, tt.setupMocks)

			value, err := sc.SelectedSetting(t.Context(), tt.key)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectValue, value)
		})
	}
}

func setupSettingsComponent(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, c *client.ClientMock, s *storage.StorageMock),
) btt.SettingsComponent {
	mc := minimock.NewController(t)
	c := client.NewClientMock(mc)
	s := storage.NewStorageMock(mc)

	if setupMocks != nil {
		setupMocks(mc, c, s)
	}

	return btt.NewSettingsComponent(c, s)
}
