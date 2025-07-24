package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFloatingMenu(t *testing.T) {
	t.Parallel()

	floatingMenu := NewFloatingMenu("foo")
	require.Equal(t, "foo", floatingMenu["BTTMenuName"])
	require.Equal(t, "foo", floatingMenu["BTTMenuConfig"].(map[string]any)["BTTMenuElementIdentifier"])
	require.Equal(t, float64(767), floatingMenu["BTTTriggerType"])
	require.Equal(t, "BTTTriggerTypeFloatingMenu", floatingMenu["BTTTriggerClass"])
}

func TestNewWebView(t *testing.T) {
	t.Parallel()

	webView := NewWebView("foo", "bar")
	require.Equal(t, "foo", webView["BTTMenuName"])
	require.Equal(t, "foo", webView["BTTMenuConfig"].(map[string]any)["BTTMenuElementIdentifier"])
	require.Equal(t, "bar", webView["BTTMenuConfig"].(map[string]any)["BTTMenuItemText"])
	require.Equal(t, float64(778), webView["BTTTriggerType"])
	require.Equal(t, "BTTTriggerTypeFloatingMenu", webView["BTTTriggerClass"])
}
