package btt_test

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"live2text/internal/services/btt"
	"live2text/internal/services/btt/client"
	"live2text/internal/services/btt/client/trigger"
	"live2text/internal/services/btt/tmpl"
)

func TestInitialize(t *testing.T) {
	t.Parallel()

	ic := setupInitializingComponent(t, func(mc *minimock.Controller, c *client.ClientMock, r *tmpl.RendererMock) {
		r.PrintStatusMock.Return("sample")
		r.PrintMetricMock.Return("sample")
		r.PrintSelectedDeviceMock.Return("sample")
		r.PrintSelectedLanguageMock.Return("sample")
		r.PrintSelectedViewModeMock.Return("sample")
		r.PrintSelectedFloatingMock.Return("sample")
		r.PrintSelectedClipboardMock.Return("sample")
		// r.SelectDeviceMock.Return("sample")
		r.SelectLanguageMock.Return("sample")
		r.SelectViewModeMock.Return("sample")
		r.SelectFloatingMock.Return("sample")
		r.FloatingPageMock.Return("sample")
		r.SelectClipboardMock.Return("sample")
		r.CloseSettingsMock.Return("sample")
		r.OpenSettingsMock.Return("sample")
		r.ToggleMock.Return("sample")
		// r.ListenSocketMock.Return("sample")
		r.AppPlaceholderMock.Return("sample")
		r.CopyTextMock.Return("sample")

		c.AddTriggerMock.Return("", nil)
	}, []string{"en"})

	err := ic.Initialize(t.Context())
	require.NoError(t, err)
}

func TestClear(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *client.ClientMock, r *tmpl.RendererMock)
		expectErr  string
	}{
		{
			name: "cannot close directory",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, r *tmpl.RendererMock) {
				c.TriggerActionMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot close directory",
		},
		{
			name: "cannot get triggers",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, r *tmpl.RendererMock) {
				c.TriggerActionMock.Return(nil)
				c.GetTriggersMock.Return(nil, errors.New("something happened"))
			},
			expectErr: "cannot get triggers",
		},
		{
			name: "cannot delete triggers",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, r *tmpl.RendererMock) {
				c.TriggerActionMock.Return(nil)
				c.GetTriggersMock.Return([]trigger.Trigger{}, nil)
				c.DeleteTriggersMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot delete triggers",
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock, r *tmpl.RendererMock) {
				c.TriggerActionMock.Return(nil)
				triggers := []trigger.Trigger{
					trigger.NewTrigger(),
				}
				c.GetTriggersMock.Expect(minimock.AnyContext, "").Return(triggers, nil)
				c.DeleteTriggersMock.Expect(minimock.AnyContext, triggers).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ic := setupInitializingComponent(t, tt.setupMocks, []string{})

			err := ic.Clear(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func setupInitializingComponent(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, c *client.ClientMock, r *tmpl.RendererMock),
	languages []string,
) btt.InitializingComponent {
	t.Helper()

	mc := minimock.NewController(t)
	c := client.NewClientMock(mc)
	r := tmpl.NewRendererMock(mc)

	if setupMocks != nil {
		setupMocks(mc, c, r)
	}

	return btt.NewInitializingComponent(c, r, languages)
}
