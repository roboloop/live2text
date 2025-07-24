package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddCycledScript(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		fn   func(script string, interval float64) Trigger
	}{
		{
			name: "internal method",
			fn:   newTrigger().addCycledScript,
		},
		{
			name: "exported method",
			fn:   newTrigger().AddCycledScript,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.fn("foo", 42)

			require.Equal(t, "/bin/bash:::-c:::-:::", trigger["BTTShellScriptWidgetGestureConfig"])
			require.Equal(
				t,
				1,
				trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarAppleScriptStringRunOnInit"],
			)
			require.Equal(t, "foo", trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarShellScriptString"])
			require.Equal(
				t,
				float64(42),
				trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarScriptUpdateInterval"],
			)

			tt.fn("foo", 0.0)
			require.Equal(
				t,
				0,
				trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarAppleScriptStringRunOnInit"],
			)
		})
	}
}

func TestAddTapScript(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		fn   func(script string) Trigger
	}{
		{
			name: "internal method",
			fn:   newTrigger().addTapScript,
		},
		{
			name: "exported method",
			fn:   newTrigger().AddTapScript,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.fn("foo")
			require.Equal(
				t,
				float64(206),
				trigger["BTTAdditionalActions"].([]any)[0].(map[string]any)["BTTPredefinedActionType"],
			)
			require.Equal(
				t,
				"foo",
				trigger["BTTAdditionalActions"].([]any)[0].(map[string]any)["BTTShellTaskActionScript"],
			)
			require.Equal(
				t,
				"/bin/bash:::-c:::-:::",
				trigger["BTTAdditionalActions"].([]any)[0].(map[string]any)["BTTShellTaskActionConfig"],
			)
		})
	}
}

func TestAddLongTapScript(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		fn   func(triggerTitle Title) Trigger
	}{
		{
			name: "internal method",
			fn:   newTrigger().addLongTapScript,
		},
		{
			name: "exported method",
			fn:   newTrigger().AddLongTapTrigger,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.fn("foo")
			require.Equal(
				t,
				"foo",
				trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarLongPressActionName"],
			)
		})
	}
}

func TestHasTapScript(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func() Trigger
		expect bool
	}{
		{
			name: "has no key",
			fn: func() Trigger {
				trigger := newTrigger()
				return trigger
			},
			expect: false,
		},
		{
			name: "has tap actions",
			fn: func() Trigger {
				trigger := newTrigger()
				trigger["BTTAdditionalActions"] = []any{}
				return trigger
			},
			expect: false,
		},
		{
			name: "has no tap action",
			fn: func() Trigger {
				trigger := newTrigger()
				trigger["BTTAdditionalActions"] = []any{
					map[string]any{
						"foo": "bar",
					},
				}
				return trigger
			},
			expect: false,
		},
		{
			name: "has tap action",
			fn: func() Trigger {
				trigger := newTrigger()
				trigger["BTTAdditionalActions"] = []any{
					map[string]any{
						"BTTPredefinedActionType": float64(206),
					},
				}
				return trigger
			},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trigger := tt.fn()
			require.Equal(t, tt.expect, trigger.HasTapScript())
		})
	}
}
