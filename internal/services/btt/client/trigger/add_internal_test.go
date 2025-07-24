package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddExtra(t *testing.T) {
	t.Parallel()

	initial := map[string]any{
		"key1": "value1",
		"key2": []string{"value2"},
		"key3": map[string]any{
			"key4": "value4",
		},
	}
	changed := map[string]any{
		"key1": 42,
		"key2": []string{"value3", "value4"},
		"key3": map[string]any{
			"key5": "value5",
		},
		"key6": "value6",
	}
	expected := map[string]any{
		"key1": 42,
		"key2": []string{"value3", "value4"},
		"key3": map[string]any{
			"key4": "value4",
			"key5": "value5",
		},
		"key6": "value6",
	}

	trigger := newTrigger()
	trigger.addExtra(initial)
	trigger.addExtra(changed)
	require.EqualValues(t, expected, trigger)
}
