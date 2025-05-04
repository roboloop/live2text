package validation_test

import (
	"live2text/internal/api/validation"
	"testing"
)

func TestIsValidLanguageCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"valid 1", "en-US", true},
		{"valid 2", "ru-RU", true},

		{"invalid 1", "invalid", false},
		{"invalid 2", "eng-US", false},
		{"invalid 3", "en-USA", false},
		{"invalid 4", "eng-USA", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validation.IsValidLanguageCode(tt.code); got != tt.want {
				t.Errorf("IsValidLanguageCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
