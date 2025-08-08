package config_test

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/config"
)

func TestParseServe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		expectCfg *config.Serve
		expectErr string
	}{
		{
			name: "default args",
			args: []string{},
			expectCfg: &config.Serve{
				AppAddress:    "127.0.0.1:8080",
				PprofAddress:  "127.0.0.1:8081",
				BttAddress:    "127.0.0.1:44444",
				OutputDir:     filepath.Join(os.TempDir(), "Live2Text"), //nolint:usetesting // business logic path
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

				"--output-dir",
				"/tmp/foo-dir",
				"--log-level",
				"debug",
				"--on-console",
			},
			expectCfg: &config.Serve{
				AppAddress:    "192.168.0.1:1234",
				PprofAddress:  "192.168.0.1:5678",
				BttAddress:    "192.168.1.1:9012",
				OutputDir:     "/tmp/foo-dir",
				LogLevel:      slog.LevelDebug,
				ConsoleWriter: os.Stdout,
			},
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

			cfg, err := config.ParseServe(io.Discard, tt.args)
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

func TestServeLogValue(t *testing.T) {
	t.Parallel()

	cfg := &config.Serve{
		AppAddress:    "192.168.0.1:1234",
		PprofAddress:  "192.168.0.1:5678",
		BttAddress:    "192.168.1.1:9012",
		OutputDir:     "/tmp/foo-dir",
		LogLevel:      slog.LevelDebug,
		ConsoleWriter: os.Stdout,
	}

	val := cfg.LogValue()
	require.Equal(t, slog.KindGroup, val.Kind())

	group := val.Group()
	expected := []slog.Attr{
		slog.String("app-address", "192.168.0.1:1234"),
		slog.String("pprof-address", "192.168.0.1:5678"),
		slog.String("btt-address", "192.168.1.1:9012"),
		slog.String("output-dir", "/tmp/foo-dir"),
		slog.String("log-level", slog.LevelDebug.String()),
		slog.String("on-console", "true"),
	}

	require.Equal(t, expected, group)
}
