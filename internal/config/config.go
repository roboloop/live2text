package config

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

type Config struct {
	AppAddress string
	BttAddress string
	Languages  []string
	LogLevel   slog.Level
}

func Initialize(args []string) (*Config, error) {
	var (
		appHost   string
		appPort   string
		bttHost   string
		bttPort   string
		languages string
		logLevel  string
	)

	fs := flag.NewFlagSet("live2text", flag.ContinueOnError)
	fs.StringVar(&appHost, "app-host", "127.0.0.1", "App server host")
	fs.StringVar(&appPort, "app-port", "8080", "App server host")
	fs.StringVar(&bttHost, "btt-host", "127.0.0.1", "BetterTouchTool webserver host")
	fs.StringVar(&bttPort, "btt-port", "44444", "BetterTouchTool webserver port")
	fs.StringVar(&languages, "languages", "en-US,es-ES,fr-FR,pt-BR,ru-RU,ja-JP,de-DE", "List of available languages")
	fs.StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn, error")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("cannot parse agruments: %w", err)
	}

	return &Config{
		net.JoinHostPort(appHost, appPort),
		net.JoinHostPort(bttHost, bttPort),
		strings.Split(languages, ","),
		parseLogLevel(logLevel),
	}, nil
}

func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
