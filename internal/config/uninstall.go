package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
)

type Uninstall struct {
	BttAddress string
	LogLevel   slog.Level
}

func (u *Uninstall) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("btt-address", u.BttAddress),
		slog.String("log-level", u.LogLevel.String()),
	)
}

func ParseUninstall(output io.Writer, args []string) (*Uninstall, error) {
	var (
		bttHost  string
		bttPort  string
		logLevel string
	)

	fs := flag.NewFlagSet("uninstall", flag.ContinueOnError)
	fs.SetOutput(output)
	fs.StringVar(&bttHost, "btt-host", "127.0.0.1", "BetterTouchTool webserver host")
	fs.StringVar(&bttPort, "btt-port", "44444", "BetterTouchTool webserver port")
	fs.StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn, error")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, ErrHelp
		}

		return nil, fmt.Errorf("cannot parse arguments: %w", err)
	}

	return &Uninstall{
		BttAddress: net.JoinHostPort(bttHost, bttPort),
		LogLevel:   parseLogLevel(logLevel),
	}, nil
}
