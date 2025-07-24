package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActionTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    actionType
		expected string
	}{
		{
			name:     "Empty Placeholder",
			input:    actionTypeEmptyPlaceholder,
			expected: "Empty Placeholder",
		},
		{
			name:     "Execute Script",
			input:    actionTypeExecuteScript,
			expected: "Execute Script",
		},
		{
			name:     "Close Group",
			input:    actionTypeCloseGroup,
			expected: "Close Group",
		},
		{
			name:     "Open Group",
			input:    actionTypeOpenGroup,
			expected: "Open Group",
		},
		{
			name:     "Open Floating Menu",
			input:    actionTypeOpenFloatingMenu,
			expected: "Open Floating Menu",
		},
		{
			name:     "Close Floating Menu",
			input:    actionTypeCloseFloatingMenu,
			expected: "Close Floating Menu",
		},
		{
			name:     "Toggle Floating Menu",
			input:    actionTypeToggleFloatingMenu,
			expected: "Toggle Floating Menu",
		},
		{
			name:     "Unknown Action Type",
			input:    actionType(9999.99),
			expected: "Unknown Action Type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.input.String()
			require.Equal(t, tt.expected, result)
		})
	}
}
