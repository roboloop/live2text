package config_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/config"
)

func TestParseUninstall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		expectCfg *config.Uninstall
		expectErr string
	}{
		{
			name: "default args",
			args: []string{},
			expectCfg: &config.Uninstall{
				BttAddress: "127.0.0.1:44444",
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
				"--btt-host",
				"192.168.1.1",
				"--btt-port",
				"9012",

				"--log-level",
				"debug",
			},
			expectCfg: &config.Uninstall{
				BttAddress: "192.168.1.1:9012",
				LogLevel:   slog.LevelDebug,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg, err := config.ParseUninstall(io.Discard, tt.args)
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

func TestUninstallLogValue(t *testing.T) {
	t.Parallel()

	cfg := config.Uninstall{
		BttAddress: "192.168.1.1:9012",
		LogLevel:   slog.LevelDebug,
	}

	val := cfg.LogValue()
	require.Equal(t, slog.KindGroup, val.Kind())

	group := val.Group()
	expected := []slog.Attr{
		slog.String("btt-address", "192.168.1.1:9012"),
		slog.String("log-level", slog.LevelDebug.String()),
	}

	require.Equal(t, expected, group)
}
