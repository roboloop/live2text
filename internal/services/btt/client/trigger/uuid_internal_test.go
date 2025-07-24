package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddUUID(t *testing.T) {
	t.Parallel()

	trigger := newTrigger().AddUUID("foo")
	require.Equal(t, "foo", trigger["BTTUUID"])
}

func TestUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupTrigger func() Trigger
		expected     UUID
	}{
		{
			name: "no value",
			setupTrigger: func() Trigger {
				trigger := newTrigger()
				return trigger
			},
			expected: UUID(""),
		},
		{
			name: "title is string",
			setupTrigger: func() Trigger {
				trigger := newTrigger()
				trigger["BTTUUID"] = "foo"
				return trigger
			},
			expected: UUID("foo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.setupTrigger()
			require.Equal(t, tt.expected, trigger.UUID())
		})
	}
}
