package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
)

type Install struct {
	AppAddress string
	BttAddress string
	LogLevel   slog.Level

	Languages []string
}

func (i *Install) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("app-address", i.AppAddress),
		slog.String("btt-address", i.BttAddress),
		slog.String("log-level", i.LogLevel.String()),
		slog.String("languages", strings.Join(i.Languages, ",")),
	)
}

func ParseInstall(output io.Writer, args []string) (*Install, error) {
	var (
		appHost string
		appPort string
		bttHost string
		bttPort string

		languages string
		logLevel  string
	)

	fs := flag.NewFlagSet("btt", flag.ContinueOnError)
	fs.SetOutput(output)
	fs.StringVar(&appHost, "app-host", "127.0.0.1", "App server host")
	fs.StringVar(&appPort, "app-port", "8080", "App server port")
	fs.StringVar(&bttHost, "btt-host", "127.0.0.1", "BetterTouchTool webserver host")
	fs.StringVar(&bttPort, "btt-port", "44444", "BetterTouchTool webserver port")
	fs.StringVar(&languages, "languages", "en-US,es-ES,fr-FR,pt-BR,ru-RU,ja-JP,de-DE", "List of available languages")
	fs.StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn, error")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, ErrHelp
		}

		return nil, fmt.Errorf("cannot parse arguments: %w", err)
	}

	return &Install{
		AppAddress: net.JoinHostPort(appHost, appPort),
		BttAddress: net.JoinHostPort(bttHost, bttPort),
		LogLevel:   parseLogLevel(logLevel),

		Languages: strings.Split(languages, ","),
	}, nil
}
