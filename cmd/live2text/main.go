// Package main is the entry point for the live2text application.
// It sets up and coordinates all the necessary services and servers.
package main

import (
	"context"
	"errors"
	"fmt"
	"live2text/internal/api"
	"live2text/internal/background"
	"live2text/internal/config"
	"live2text/internal/services"
	"live2text/internal/services/audio"
	"live2text/internal/services/audio_wrapper"
	"live2text/internal/services/btt"
	bttexec "live2text/internal/services/btt/exec"
	btthttp "live2text/internal/services/btt/http"
	"live2text/internal/services/burner"
	"live2text/internal/services/metrics"
	"live2text/internal/services/recognition"
	"live2text/internal/services/speech_wrapper"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args[1:]); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Initialize(args)
	if err != nil {
		return fmt.Errorf("cannot initialize config: %w", err)
	}

	sc, err := speech_wrapper.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("cannot create speech client: %w", err)
	}
	defer sc.Close()

	aw, awClose, err := audio_wrapper.NewAudio()
	if err != nil {
		return fmt.Errorf("cannot create new audio wrapper: %w", err)
	}
	defer awClose()

	sl, l := newLoggers(cfg.LogLevel)

	tm := background.NewTaskManager(ctx)
	sm := background.NewSocketManager(ctx, sl)

	m := metrics.NewMetrics()
	a := audio.NewAudio(sl, m, aw)
	b := burner.NewBurner(sl, m)
	r := recognition.NewRecognition(sl, m, a, b, sc, tm, sm)
	b2 := newBtt(sl, a, r, cfg)

	s := services.NewServices(a, aw, b, r, m, b2)

	server := newServer(ctx, cfg.AppAddress, api.NewHandler(sl, s), l)
	defer server.Close()
	go func() {
		slog.InfoContext(ctx, "Starting the API server", "address", cfg.AppAddress)
		if serverErr := server.ListenAndServe(); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "Failed to listen and serve", "error", serverErr)
		}
	}()

	debugServer := newServer(ctx, "127.0.0.1:6060", newDebugHandler(), l)
	go func() {
		slog.InfoContext(ctx, "Starting debug server", "address", debugServer.Addr)
		if debugServerErr := debugServer.ListenAndServe(); debugServerErr != nil && !errors.Is(debugServerErr, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "Debug server error", "error", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	slog.InfoContext(ctx, "Received shutdown signal", "signal", ctx.Err())

	// Shutdown servers and clean up resource
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown server gracefully", "error", err)
	} else {
		slog.InfoContext(ctx, "Server shutdown completed successfully")
	}

	if err = debugServer.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown debug server gracefully", "error", err)
	} else {
		slog.InfoContext(ctx, "Debug server shutdown completed successfully")
	}

	slog.InfoContext(ctx, "Closing socket manager")
	sm.Close()

	slog.InfoContext(ctx, "Waiting for background tasks to complete")
	tm.Wait()

	slog.InfoContext(ctx, "Application shutdown complete")

	return nil
}

func newLoggers(level slog.Level) (*slog.Logger, *log.Logger) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler), slog.NewLogLogger(handler, level)
}

func newBtt(logger *slog.Logger, audio audio.Audio, recognition recognition.Recognition, cfg *config.Config) btt.Btt {
	const (
		appName = "live2text"
		bttName = "BetterTouchTool"
	)

	httpClient := btthttp.NewClient(logger, cfg.BttAddress)
	execClient := bttexec.NewClient(logger, appName, bttName)

	return btt.NewBtt(logger, audio, recognition, httpClient, execClient, cfg)
}

func newServer(ctx context.Context, address string, handler http.Handler, logger *log.Logger) *http.Server {
	return &http.Server{
		Addr:         address,
		Handler:      handler,
		ErrorLog:     logger,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}
}

func newDebugHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	return mux
}
