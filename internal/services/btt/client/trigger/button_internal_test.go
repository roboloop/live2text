package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTapButton(t *testing.T) {
	t.Parallel()

	button := NewTapButton("foo", "bar")

	requireTitle(t, "foo", button)
	require.Equal(t, float64(629), button["BTTTriggerType"])
	require.Equal(t, "BTTTriggerTypeTouchBar", button["BTTTriggerClass"])
	require.Equal(t, float64(366), button["BTTPredefinedActionType"])
	require.Equal(t, "bar", button["BTTAdditionalActions"].([]any)[0].(map[string]any)["BTTShellTaskActionScript"])
}

func TestNewTapIconButton(t *testing.T) {
	t.Parallel()

	button := NewTapIconButton("foo", "bar", "baz")
	requireTitle(t, "foo", button)

	require.Equal(t, "baz", button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemSFSymbolDefaultIcon"])
	require.Equal(t, float64(22), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemIconHeight"])
	require.Equal(t, true, button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarOnlyShowIcon"])
}

func TestNewInfoButton(t *testing.T) {
	t.Parallel()

	button := NewInfoButton("foo", "bar", 42)

	requireTitle(t, "foo", button)
	require.Equal(t, float64(642), button["BTTTriggerType"])
	require.Equal(t, "BTTTriggerTypeTouchBar", button["BTTTriggerClass"])
	require.Equal(t, float64(366), button["BTTPredefinedActionType"])
	require.Equal(t, "bar", button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarShellScriptString"])
	require.Equal(t, float64(42), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarScriptUpdateInterval"])
}

func TestNewStatusInfoButton(t *testing.T) {
	t.Parallel()

	button := NewStatusInfoButton("foo", "bar")

	requireTitle(t, "foo", button)
	require.Equal(t, "bar", button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarShellScriptString"])
	require.Equal(t, float64(15), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarScriptUpdateInterval"])
	require.NotContains(t, button["BTTTriggerConfig"], "BTTTouchBarFreeSpaceAfterButton")
}

func TestNewSettingsInfoButton(t *testing.T) {
	t.Parallel()

	button := NewSettingsInfoButton("foo", "bar")

	requireTitle(t, "foo", button)
	require.Equal(t, "bar", button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarShellScriptString"])
	require.Equal(t, float64(15), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarScriptUpdateInterval"])
	require.Equal(t, float64(25), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarFreeSpaceAfterButton"])
}

func TestNewMetricsInfoButton(t *testing.T) {
	t.Parallel()

	button := NewMetricsInfoButton("foo", "bar", "baz")

	requireTitle(t, "foo", button)
	require.Equal(t, "bar", button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarShellScriptString"])
	require.Equal(t, float64(5), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarScriptUpdateInterval"])

	require.Equal(t, "baz", button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemSFSymbolDefaultIcon"])
	require.Equal(t, float64(22), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemIconHeight"])
	require.NotContains(t, button["BTTTriggerConfig"], "BTTTouchBarOnlyShowIcon")
}

func TestNewHiddenDir(t *testing.T) {
	t.Parallel()

	button := NewHiddenDir("foo")

	requireTitle(t, "foo", button)
	require.Equal(t, float64(630), button["BTTTriggerType"])
	require.Equal(t, "BTTTriggerTypeTouchBar", button["BTTTriggerClass"])
	require.Equal(t, 0, button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarButtonWidth"])
	require.Equal(t, 1, button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarButtonUseFixedWidth"])
}

func TestNewDirButton(t *testing.T) {
	t.Parallel()

	button := NewDirButton("foo", "bar")

	requireTitle(t, "foo", button)
	require.Equal(t, float64(630), button["BTTTriggerType"])
	require.Equal(t, "BTTTriggerTypeTouchBar", button["BTTTriggerClass"])
	require.Equal(t, "bar", button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemSFSymbolDefaultIcon"])
	require.Equal(t, float64(22), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemIconHeight"])
	require.NotContains(t, button["BTTTriggerConfig"], "BTTTouchBarOnlyShowIcon")
}

func TestNewCloseDirButton(t *testing.T) {
	t.Parallel()

	button := NewCloseDirButton()

	requireTitle(t, "Close Directory", button)
	require.Equal(t, float64(629), button["BTTTriggerType"])
	require.Equal(t, "BTTTriggerTypeTouchBar", button["BTTTriggerClass"])
	require.Equal(t, float64(191), button["BTTPredefinedActionType"])

	require.Equal(
		t,
		"xmark.circle.fill",
		button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemSFSymbolDefaultIcon"],
	)
	require.Equal(t, float64(25), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemIconHeight"])
	require.Equal(t, true, button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarOnlyShowIcon"])
}

func TestNewOpenDirButton(t *testing.T) {
	t.Parallel()

	button := NewOpenDirButton("foo")

	requireTitle(t, "Close Directory", button)
	require.Equal(t, float64(629), button["BTTTriggerType"])
	require.Equal(t, "BTTTriggerTypeTouchBar", button["BTTTriggerClass"])
	require.Equal(t, float64(205), button["BTTPredefinedActionType"])

	require.Equal(t, "foo", button["BTTOpenGroupWithName"])
	require.Equal(
		t,
		"xmark.circle.fill",
		button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemSFSymbolDefaultIcon"],
	)
	require.Equal(t, float64(25), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemIconHeight"])
	require.Equal(t, true, button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarOnlyShowIcon"])
}

func TestNewNamedTrigger(t *testing.T) {
	t.Parallel()

	namedTrigger := NewNamedTrigger("foo", "bar")

	requireTitle(t, "foo", namedTrigger)
	require.Equal(t, float64(643), namedTrigger["BTTTriggerType"])
	require.Equal(t, "BTTTriggerTypeOtherTriggers", namedTrigger["BTTTriggerClass"])
	require.Equal(t, float64(366), namedTrigger["BTTPredefinedActionType"])
	require.Equal(
		t,
		"bar",
		namedTrigger["BTTAdditionalActions"].([]any)[0].(map[string]any)["BTTShellTaskActionScript"],
	)
}

func TestAddCloseIcon(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		fn   func(Trigger) Trigger
	}{
		{
			name: "internal method",
			fn:   Trigger.addCloseIcon,
		},
		{
			name: "exported method",
			fn:   Trigger.AddCloseIcon,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			button := tt.fn(newTrigger())
			require.Equal(
				t,
				"xmark.circle.fill",
				button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemSFSymbolDefaultIcon"],
			)
			require.Equal(t, float64(25), button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarItemIconHeight"])
			require.Equal(t, true, button["BTTTriggerConfig"].(map[string]any)["BTTTouchBarOnlyShowIcon"])
		})
	}
}

func requireTitle(t *testing.T, expected string, button Trigger) {
	t.Helper()

	require.Equal(t, expected, button["BTTTouchBarButtonName"])
	require.Equal(t, expected, button["BTTWidgetName"])
	require.Equal(t, expected, button["BTTTriggerName"])
	require.Equal(t, expected, button["BTTMenuName"])
}
