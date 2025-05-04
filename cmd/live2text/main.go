package main

import (
	"context"
	"fmt"
	"live2text/internal/api"
	"live2text/internal/background"
	"live2text/internal/config"
	"live2text/internal/services"
	"live2text/internal/services/audio"
	"live2text/internal/services/audio_wrapper"
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
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
		server := http.Server{
			Addr:        "127.0.0.1:6060",
			Handler:     mux,
			ReadTimeout: 5 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Debug server error", "error", err)
		}
	}()

	ctx := context.Background()
	if err := run(ctx, os.Args[1:]); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Initialize(args)
	if err != nil {
		return fmt.Errorf("could not initialize config: %w", err)
	}

	sc, err := speech_wrapper.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("could not create speech client: %w", err)
	}
	defer sc.Close()

	aw, awClose, err := audio_wrapper.NewAudio()
	if err != nil {
		return fmt.Errorf("cannot create new aduio wrapper: %w", err)
	}
	defer awClose()

	sl, l := newLoggers(slog.LevelInfo)
	tm := background.NewTaskManager(ctx)
	sm := background.NewSocketManager(ctx, sl)
	m := metrics.NewMetrics()
	a := audio.NewAudio(sl, m, aw)
	b := burner.NewBurner(sl, m)
	r := recognition.NewRecognition(sl, m, a, b, sc, tm, sm)
	s := services.NewServices(a, aw, b, r, m)

	server := &http.Server{
		Addr:     net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:  api.NewHandler(sl, s),
		ErrorLog: l,
	}
	defer server.Close()
	go func() {
		slog.InfoContext(ctx, "Starting the api server")
		if serverErr := server.ListenAndServe(); serverErr != nil {
			slog.ErrorContext(ctx, "Failed to listen and serve", "error", serverErr)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown server", "error", err)
	}
	sm.Close()
	tm.Wait()

	slog.InfoContext(ctx, "Shutting down the program", "error", ctx.Err())

	return nil
}

func newLoggers(level slog.Level) (*slog.Logger, *log.Logger) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler), slog.NewLogLogger(handler, level)
}
