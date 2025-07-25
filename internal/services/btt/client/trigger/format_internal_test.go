package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddReadableFormat(t *testing.T) {
	t.Parallel()

	trigger := newTrigger()
	trigger = trigger.AddReadableFormat()
	require.Equal(t, 12, trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarButtonFontSize"])
	require.Equal(
		t,
		"0.0, 0.0, 0.0, 255.0",
		trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarButtonColor"],
	)
	require.Equal(t, 0, trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarButtonTextAlignment"])
}

func TestAddEnabled(t *testing.T) {
	t.Parallel()

	patch := newTrigger().AddEnabled()
	require.Equal(t, 1, patch["BTTEnabled2"])
}

func TestAddDisabled(t *testing.T) {
	t.Parallel()

	patch := newTrigger().AddDisabled()
	require.Equal(t, 0, patch["BTTEnabled2"])
}
