package config_test

import (
	"live2text/internal/config"
	"log/slog"
	"reflect"
	"testing"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expected    *config.Config
		expectedErr string
	}{
		{
			name: "Default args",
			args: []string{},
			expected: &config.Config{
				AppAddress: "127.0.0.1:8080",
				BttAddress: "127.0.0.1:44444",
				Languages:  []string{"en-US", "es-ES", "fr-FR", "pt-BR", "ru-RU", "ja-JP", "de-DE"},
				LogLevel:   slog.LevelInfo,
			},
		},
		{
			name: "Passed args",
			args: []string{
				"--app-host",
				"192.168.0.1",
				"--app-port",
				"1234",
				"--btt-host",
				"192.168.1.1",
				"--btt-port",
				"5678",
				"--languages",
				"es-ES,pt-BR",
				"--log-level",
				"debug",
			},
			expected: &config.Config{
				AppAddress: "192.168.0.1:1234",
				BttAddress: "192.168.1.1:5678",
				Languages:  []string{"es-ES", "pt-BR"},
				LogLevel:   slog.LevelDebug,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Initialize(tt.args)
			if tt.expectedErr != "" {
				if err != nil {
					t.Fatalf("Initialize() expected error %v, got nil", tt.expectedErr)
				}

				if tt.expectedErr != err.Error() {
					t.Errorf("Expected error: %v, got %v", tt.expectedErr, err.Error())
				}
				return
			}

			if !reflect.DeepEqual(cfg, tt.expected) {
				t.Errorf("Initialize() got = %v, want %v", cfg, tt.expected)
			}
		})
	}
}
