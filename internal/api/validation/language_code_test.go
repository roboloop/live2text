package validation_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/api/validation"
)

func TestIsValidLanguageCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"valid 1", "en-US", true},
		{"valid 2", "es-ES", true},

		{"invalid 1", "invalid", false},
		{"invalid 2", "eng-US", false},
		{"invalid 3", "en-USA", false},
		{"invalid 4", "eng-USA", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := validation.IsValidLanguageCode(tt.code)
			require.Equal(t, tt.expected, result)
		})
	}
}
