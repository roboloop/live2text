package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCloseDirAction(t *testing.T) {
	t.Parallel()

	action := NewCloseDirAction()
	require.Equal(t, float64(191), action["BTTPredefinedActionType"])
}

func TestNewOpenDirAction(t *testing.T) {
	t.Parallel()

	action := NewOpenDirAction("foo")
	require.Equal(t, "foo", action["BTTOpenGroupWithName"])
}

func TestNewOpenFloatingAction(t *testing.T) {
	t.Parallel()

	action := NewOpenFloatingAction("foo")
	require.Equal(t, float64(386), action["BTTPredefinedActionType"])
	require.Equal(t, "foo", action["BTTAdditionalActionData"].(map[string]any)["BTTMenuActionMenuID"])
}

func TestNewCloseFloatingAction(t *testing.T) {
	t.Parallel()

	action := NewCloseFloatingAction("foo")
	require.Equal(t, float64(387), action["BTTPredefinedActionType"])
	require.Equal(t, "foo", action["BTTAdditionalActionData"].(map[string]any)["BTTMenuActionMenuID"])
}

func TestAddCloseDirAction(t *testing.T) {
	t.Parallel()

	trigger := newTrigger().addCloseDirAction()
	require.Equal(t, float64(191), trigger["BTTPredefinedActionType"])
}

func TestAddOpenDirAction(t *testing.T) {
	t.Parallel()

	trigger := newTrigger().addOpenDirAction("foo")
	require.Equal(t, float64(205), trigger["BTTPredefinedActionType"])
	require.Equal(t, "foo", trigger["BTTOpenGroupWithName"])
}

func TestAddAction(t *testing.T) {
	t.Parallel()

	trigger := newTrigger().addAction(42)
	require.Equal(t, float64(42), trigger["BTTPredefinedActionType"])
}
