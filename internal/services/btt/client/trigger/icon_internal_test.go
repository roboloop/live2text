package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddIcon(t *testing.T) {
	t.Parallel()

	icon := newTrigger().addIcon("foo", 42, false)
	config := icon["BTTTriggerConfig"].(map[string]any)
	require.Equal(t, "foo", config["BTTTouchBarItemSFSymbolDefaultIcon"])
	require.Equal(t, 0, config["BTTTouchBarItemSFSymbolWeight"])
	require.Equal(t, 2, config["BTTTouchBarItemIconType"])
	require.Equal(t, float64(42), config["BTTTouchBarItemIconHeight"])
	require.Equal(t, -10, config["BTTTouchBarItemPadding"])

	onlyIcon := newTrigger().addIcon("foo", 42, true)
	config = onlyIcon["BTTTriggerConfig"].(map[string]any)
	require.Equal(t, "0.0, 0.0, 0.0, 255.0", config["BTTTouchBarButtonColor"])
	require.Equal(t, true, config["BTTTouchBarOnlyShowIcon"])
}
