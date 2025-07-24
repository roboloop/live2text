package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddLabel(t *testing.T) {
	t.Parallel()

	button := newTrigger()
	button = button.AddLabel("foo")
	require.Equal(t, "foo", button["BTTGroupName"])
	require.Equal(t, "foo", button["BTTNotes"])

	other := newTrigger().addExtra(map[string]any{"BTTTriggerClass": typeOtherTriggers})
	other = other.AddLabel("bar")
	require.Equal(t, "bar", other["BTTGestureNotes"])
}
