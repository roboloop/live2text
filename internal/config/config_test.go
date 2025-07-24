package config_test

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"live2text/internal/config"
)

func TestInitialize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		expectCfg *config.Config
		expectErr string
	}{
		{
			name: "default args",
			args: []string{},
			expectCfg: &config.Config{
				AppAddress:    "127.0.0.1:8080",
				PprofAddress:  "127.0.0.1:8081",
				BttAddress:    "127.0.0.1:44444",
				AppName:       "Live2Text",
				OutputDir:     os.TempDir() + "live2text", //nolint:usetesting // business logic path
				LogLevel:      slog.LevelInfo,
				ConsoleWriter: io.Discard,
			},
		},
		{
			name: "passed args",
			args: []string{
				"--app-host",
				"192.168.0.1",
				"--app-port",
				"1234",

				"--pprof-host",
				"192.168.0.1",
				"--pprof-port",
				"5678",

				"--btt-host",
				"192.168.1.1",
				"--btt-port",
				"9012",

				"--app-name",
				"foo-name",
				"--output-dir",
				"/tmp/foo-dir",
				"--log-level",
				"debug",
				"--on-console",
			},
			expectCfg: &config.Config{
				AppAddress:    "192.168.0.1:1234",
				PprofAddress:  "192.168.0.1:5678",
				BttAddress:    "192.168.1.1:9012",
				AppName:       "foo-name",
				OutputDir:     "/tmp/foo-dir",
				LogLevel:      slog.LevelDebug,
				ConsoleWriter: os.Stdout,
			},
		},
		{
			name:      "invalid arguments",
			args:      []string{"--invalid"},
			expectCfg: nil,
			expectErr: "cannot parse arguments",
		},
		{
			name: "cannot create the output dir",
			args: []string{
				"--output-dir",
				string([]byte{0x00}),
			},
			expectErr: "cannot create the output dir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.name == "default args" {
				t.Skip("doesn't work in CI")
			}

			cfg, err := config.Initialize(io.Discard, tt.args)
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

func TestInitializeBtt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		expectCfg *config.BttConfig
		expectErr string
	}{
		{
			name: "default args",
			args: []string{},
			expectCfg: &config.BttConfig{
				AppAddress: "127.0.0.1:8080",
				BttAddress: "127.0.0.1:44444",
				Languages:  []string{"en-US", "es-ES", "fr-FR", "pt-BR", "ru-RU", "ja-JP", "de-DE"},
				AppName:    "Live2Text",
				LogLevel:   slog.LevelInfo,
				Clear:      false,
			},
			expectErr: "",
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
				"--app-name",
				"foo-name",
				"--log-level",
				"debug",
				"--clear",
			},
			expectCfg: &config.BttConfig{
				AppAddress: "192.168.0.1:1234",
				BttAddress: "192.168.1.1:9012",
				Languages:  []string{"es-ES", "pt-BR"},
				AppName:    "foo-name",
				LogLevel:   slog.LevelDebug,
				Clear:      true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg, err := config.InitializeBtt(io.Discard, tt.args)
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
