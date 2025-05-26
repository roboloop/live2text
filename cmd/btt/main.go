package main

import (
	"context"
	"fmt"
	"live2text/internal/config"
	"live2text/internal/services/btt"
	"live2text/internal/services/btt/exec"
	"live2text/internal/services/btt/http"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args[1:]); err != nil {
		slog.Error("Application error", "error", err) //nolint:sloglint
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	cfg, err := config.Initialize(args)
	if err != nil {
		return fmt.Errorf("cannot initialize config: %w", err)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	})
	logger := slog.New(handler)
	bttClient := newBtt(logger, cfg)
	if err = bttClient.Clear(ctx); err != nil {
		return fmt.Errorf("cannot clear: %w", err)
	}
	if err = bttClient.Initialize(ctx); err != nil {
		return fmt.Errorf("cannot initialize: %w", err)
	}

	return nil
}

func newBtt(logger *slog.Logger, cfg *config.Config) btt.Btt {
	const (
		appName = "live2text"
		bttName = "BetterTouchTool"
	)

	httpClient := http.NewClient(logger, cfg.BttAddress)
	execClient := exec.NewClient(logger, appName, bttName)

	return btt.NewBtt(logger, nil, nil, httpClient, execClient, cfg)
}
