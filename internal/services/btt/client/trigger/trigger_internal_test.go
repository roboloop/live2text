package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTrigger(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		fn   func() Trigger
	}{
		{
			name: "internal method",
			fn:   newTrigger,
		},
		{
			name: "exported method",
			fn:   NewTrigger,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.fn()
			require.Empty(t, trigger)
		})
	}
}

func TestInit(t *testing.T) {
	t.Parallel()

	trigger := newTrigger().init("foo", 42, "bar")
	require.Equal(t, "foo", trigger["BTTTouchBarButtonName"])
	require.Equal(t, "foo", trigger["BTTWidgetName"])
	require.Equal(t, "foo", trigger["BTTTriggerName"])
	require.Equal(t, "foo", trigger["BTTMenuName"])
	require.Equal(t, float64(42), trigger["BTTTriggerType"])
	require.Equal(t, "bar", trigger["BTTTriggerClass"])
	require.Equal(t, true, trigger["BTTTriggerConfig"].(map[string]any)["BTTKeepGroupOpenWhileSwitchingApps"])
}

func TestErrorContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupTrigger func() Trigger
		expected     string
	}{
		{
			name: "by name",
			setupTrigger: func() Trigger {
				trigger := newTrigger()
				trigger["BTTTriggerName"] = "foo"
				return trigger
			},
			expected: "trigger foo",
		},
		{
			name: "by action",
			setupTrigger: func() Trigger {
				trigger := newTrigger()
				trigger["BTTPredefinedActionType"] = float64(366)
				return trigger
			},
			expected: "action Empty Placeholder",
		},
		{
			name: "unknown",
			setupTrigger: func() Trigger {
				trigger := newTrigger()
				return trigger
			},
			expected: "unknown trigger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.setupTrigger()
			require.Equal(t, tt.expected, trigger.ErrorContext())
		})
	}
}

func TestTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupTrigger func() Trigger
		expected     Title
	}{
		{
			name: "no value",
			setupTrigger: func() Trigger {
				trigger := newTrigger()
				return trigger
			},
			expected: Title(""),
		},
		{
			name: "title is string",
			setupTrigger: func() Trigger {
				trigger := newTrigger()
				trigger["BTTTriggerName"] = "foo"
				return trigger
			},
			expected: Title("foo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.setupTrigger()
			require.Equal(t, tt.expected, trigger.Title())
		})
	}
}

func TestActionType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupTrigger func() Trigger
		expect       actionType
	}{
		{
			name: "no value",
			setupTrigger: func() Trigger {
				trigger := newTrigger()
				return trigger
			},
			expect: actionType(0),
		},
		{
			name: "title is string",
			setupTrigger: func() Trigger {
				trigger := newTrigger()
				trigger["BTTPredefinedActionType"] = float64(42)
				return trigger
			},
			expect: actionType(42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.setupTrigger()
			require.Equal(t, tt.expect, trigger.actionType())
		})
	}
}
