package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"golang.org/x/term"

	"live2text/internal/config"
	"live2text/internal/env"
	"live2text/internal/services/btt"
	bttclient "live2text/internal/services/btt/client"
	httpclient "live2text/internal/services/btt/client/http"
	"live2text/internal/services/btt/tmpl"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args[1:]); err != nil {
		slog.Error("Application error", "error", err) //nolint:sloglint
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	cfg, err := config.InitializeBtt(os.Stderr, args)
	if err != nil {
		return fmt.Errorf("cannot initialize config: %w", err)
	}

	logger := newLogger(cfg.LogLevel)
	ic := newInitializingComponent(logger, cfg)
	if err = ic.Clear(ctx); err != nil {
		return fmt.Errorf("cannot clear: %w", err)
	}
	if !cfg.Clear {
		if err = ic.Initialize(ctx); err != nil {
			return fmt.Errorf("cannot initialize: %w", err)
		}
	}

	return nil
}

func newLogger(level slog.Level) *slog.Logger {
	noColor := !term.IsTerminal(int(os.Stderr.Fd()))

	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:   level,
		NoColor: noColor,
	})

	return slog.New(handler)
}

func newInitializingComponent(logger *slog.Logger, cfg *config.BttConfig) btt.InitializingComponent {
	httpClient := httpclient.NewClient(logger, cfg.BttAddress, nil)
	client := bttclient.NewClient(httpClient, cfg.AppName)

	renderer := tmpl.NewRenderer(cfg.AppName, cfg.AppAddress, cfg.BttAddress, env.IsDebugMode())

	return btt.NewInitializingComponent(client, renderer, cfg.Languages)
}
