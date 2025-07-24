package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddOrder(t *testing.T) {
	t.Parallel()

	trigger := newTrigger()
	require.NotContains(t, trigger, "BTTOrder")

	trigger.AddOrder(42)
	require.Equal(t, float64(42), trigger["BTTOrder"])
}

func TestAddOrderAfter(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name        string
		trigger     Trigger
		after       Trigger
		expectOrder float64
	}{
		{
			name:        "after is nil means nothing",
			trigger:     newTrigger(),
			after:       nil,
			expectOrder: float64(0),
		},
		{
			name:        "after is empty means 0",
			trigger:     newTrigger(),
			after:       newTrigger(),
			expectOrder: float64(1),
		},
		{
			name:        "after is not nil",
			trigger:     newTrigger(),
			after:       newTrigger().AddOrder(42),
			expectOrder: float64(43),
		},
		{
			name:        "order overwritten",
			trigger:     newTrigger().AddOrder(42),
			after:       newTrigger().AddOrder(1337),
			expectOrder: float64(1338),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.trigger.AddOrderAfter(tt.after)
			require.Equal(t, tt.expectOrder, trigger["BTTOrder"])
		})
	}
}
