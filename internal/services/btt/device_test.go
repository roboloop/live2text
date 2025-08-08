package btt_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/audio"
	"github.com/roboloop/live2text/internal/services/btt"
	"github.com/roboloop/live2text/internal/services/btt/client"
	"github.com/roboloop/live2text/internal/services/btt/client/trigger"
	"github.com/roboloop/live2text/internal/services/btt/tmpl"
)

func TestLoadDevices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock)
		expectErr  string
	}{
		{
			name: "cannot get device dir trigger",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
				c.GetTriggerMock.Return(nil, errors.New("something happened"))
			},
			expectErr: "cannot get device dir trigger",
		},
		{
			name: "cannot get device triggers",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
				c.GetTriggerMock.Return(trigger.NewTrigger(), nil)
				c.GetTriggersMock.Return(nil, errors.New("something happened"))
			},
			expectErr: "cannot get device triggers",
		},
		{
			name: "cannot delete triggers",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
				c.GetTriggerMock.Return(trigger.NewTrigger(), nil)
				c.GetTriggersMock.Return([]trigger.Trigger{}, nil)
				c.DeleteTriggersMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot delete triggers",
		},
		{
			name: "cannot get a list of devices",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
				c.GetTriggerMock.Return(trigger.NewTrigger(), nil)
				c.GetTriggersMock.Return([]trigger.Trigger{}, nil)
				c.DeleteTriggersMock.Return(nil)
				a.ListOfNamesMock.Return(nil, errors.New("something happened"))
			},
			expectErr: "cannot get a list of devices",
		},
		{
			name: "cannot create trigger",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
				c.GetTriggerMock.Return(trigger.NewTrigger(), nil)
				c.GetTriggersMock.Return([]trigger.Trigger{}, nil)
				c.DeleteTriggersMock.Return(nil)
				a.ListOfNamesMock.Return([]string{"foo"}, nil)
				r.SelectDeviceMock.Return("baz")
				c.AddTriggerMock.Return("", errors.New("something happened"))
			},
			expectErr: "cannot create trigger",
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
				triggers := []trigger.Trigger{
					trigger.NewCloseDirButton(),
					trigger.NewInfoButton("foo", "bar", 42),
					trigger.NewTapButton("baz", "qux"),
				}
				deviceDir := trigger.NewTrigger().AddUUID("quux")
				selectedDevice := trigger.NewTrigger().AddOrder(1)

				getTriggerCalls := 0
				c.GetTriggerMock.Set(func(ctx context.Context, title trigger.Title) (trigger.Trigger, error) {
					getTriggerCalls += 1
					switch getTriggerCalls {
					case 1:
						require.Equal(t, "Device", title.String())
						return deviceDir, nil
					default:
						require.Equal(t, "Selected Device", title.String())
						return selectedDevice, nil
					}
				})
				c.GetTriggersMock.Expect(minimock.AnyContext, "quux").Return(triggers, nil)
				c.DeleteTriggersMock.Expect(minimock.AnyContext, triggers[2:]).Return(nil)
				a.ListOfNamesMock.Return([]string{"foo"}, nil)
				r.SelectDeviceMock.Expect("foo").Return("baz")
				c.AddTriggerMock.Inspect(func(ctx context.Context, trigger trigger.Trigger, parentUUID trigger.UUID) {
					require.Equal(t, "foo", trigger.Title().String())
					require.Equal(t, "quux", parentUUID.String())
				}).Return("00000000-0000-0000-0000-000000000000", nil)
			},
			expectErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dc := setupDeviceComponent(t, tt.setupMocks)

			err := dc.LoadDevices(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestSelectDevice(t *testing.T) {
	t.Parallel()

	dc := setupDeviceComponent(
		t,
		func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
			s.SelectSettingsMock.Expect(minimock.AnyContext, "Selected Device", "LIVE2TEXT_SELECTED_DEVICE", "foo").
				Return(nil)
		},
	)
	err := dc.SelectDevice(t.Context(), "foo")
	require.NoError(t, err)
}

func TestSelectedDevice(t *testing.T) {
	t.Parallel()

	dc := setupDeviceComponent(
		t,
		func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
			s.SelectedSettingMock.Expect(minimock.AnyContext, "LIVE2TEXT_SELECTED_DEVICE").Return("foo", nil)
		},
	)
	device, err := dc.SelectedDevice(t.Context())
	require.NoError(t, err)
	require.Equal(t, "foo", device)
}

func TestIsAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		setupMocks      func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock)
		device          string
		expectAvailable bool
		expectErr       string
	}{
		{
			name: "cannot get a list of devices",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
				a.ListOfNamesMock.Return(nil, errors.New("something happened"))
			},
			expectErr: "cannot get a list of devices",
		},
		{
			name: "device is unavailable",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
				a.ListOfNamesMock.Return([]string{"foo", "bar"}, nil)
			},
			device:          "baz",
			expectAvailable: false,
		},
		{
			name: "device is available",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock) {
				a.ListOfNamesMock.Return([]string{"foo", "bar"}, nil)
			},
			device:          "foo",
			expectAvailable: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dc := setupDeviceComponent(t, tt.setupMocks)
			available, err := dc.IsAvailable(t.Context(), tt.device)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectAvailable, available)
		})
	}
}

func setupDeviceComponent(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, a *audio.AudioMock, c *client.ClientMock, r *tmpl.RendererMock, s *btt.SettingsComponentMock),
) btt.DeviceComponent {
	mc := minimock.NewController(t)
	a := audio.NewAudioMock(mc)
	c := client.NewClientMock(mc)
	r := tmpl.NewRendererMock(mc)
	s := btt.NewSettingsComponentMock(mc)
	if setupMocks != nil {
		setupMocks(mc, a, c, r, s)
	}

	return btt.NewDeviceComponent(a, c, r, s)
}
