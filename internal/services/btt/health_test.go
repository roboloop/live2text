package btt_test

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/btt"
	"github.com/roboloop/live2text/internal/services/btt/client"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupMocks   func(mc *minimock.Controller, c *client.ClientMock)
		expectHealth bool
	}{
		{
			name: "check health",
			setupMocks: func(mc *minimock.Controller, c *client.ClientMock) {
				c.HealthMock.Expect(minimock.AnyContext).Return(true)
			},
			expectHealth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hc := setupHealthComponent(t, tt.setupMocks)

			health := hc.Health(t.Context())
			require.True(t, health)
		})
	}
}

func setupHealthComponent(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, c *client.ClientMock),
) btt.HealthComponent {
	mc := minimock.NewController(t)
	c := client.NewClientMock(mc)

	if setupMocks != nil {
		setupMocks(mc, c)
	}

	return btt.NewHealthComponent(c)
}
