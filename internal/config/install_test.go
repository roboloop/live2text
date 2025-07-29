package config_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"live2text/internal/config"
)

func TestParseInstall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		expectCfg *config.Install
		expectErr string
	}{
		{
			name: "default args",
			args: []string{},
			expectCfg: &config.Install{
				AppAddress: "127.0.0.1:8080",
				BttAddress: "127.0.0.1:44444",
				Languages:  []string{"en-US", "es-ES", "fr-FR", "pt-BR", "ru-RU", "ja-JP", "de-DE"},
				LogLevel:   slog.LevelInfo,
			},
			expectErr: "",
		},
		{
			name:      "help shown",
			args:      []string{"--help"},
			expectErr: "help shown",
		},
		{
			name:      "invalid arguments",
			args:      []string{"--invalid"},
			expectCfg: nil,
			expectErr: "cannot parse arguments",
		},
		{
			name: "passed args",
			args: []string{
				"--app-host",
				"192.168.0.1",
				"--app-port",
				"1234",

				"--btt-host",
				"192.168.1.1",
				"--btt-port",
				"9012",

				"--languages",
				"es-ES,pt-BR",
				"--log-level",
				"debug",
			},
			expectCfg: &config.Install{
				AppAddress: "192.168.0.1:1234",
				BttAddress: "192.168.1.1:9012",
				Languages:  []string{"es-ES", "pt-BR"},
				LogLevel:   slog.LevelDebug,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg, err := config.ParseInstall(io.Discard, tt.args)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectCfg, cfg)
		})
	}
}

func TestInstallLogValue(t *testing.T) {
	t.Parallel()

	cfg := config.Install{
		AppAddress: "192.168.0.1:1234",
		BttAddress: "192.168.1.1:9012",
		LogLevel:   slog.LevelDebug,
		Languages:  []string{"es-ES", "pt-BR"},
	}

	val := cfg.LogValue()
	require.Equal(t, slog.KindGroup, val.Kind())

	group := val.Group()
	expected := []slog.Attr{
		slog.String("app-address", "192.168.0.1:1234"),
		slog.String("btt-address", "192.168.1.1:9012"),
		slog.String("log-level", slog.LevelDebug.String()),
		slog.String("languages", "es-ES,pt-BR"),
	}

	require.Equal(t, expected, group)
}
