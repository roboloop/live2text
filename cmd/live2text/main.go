// Package main is the entry point for the live2text application.
// It sets up and coordinates all the necessary services and servers.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
	"golang.org/x/term"

	"live2text/internal/api"
	"live2text/internal/background"
	"live2text/internal/config"
	"live2text/internal/env"
	"live2text/internal/services"
	"live2text/internal/services/audio"
	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/services/btt"
	bttclient "live2text/internal/services/btt/client"
	httpclient "live2text/internal/services/btt/client/http"
	"live2text/internal/services/btt/storage"
	"live2text/internal/services/btt/tmpl"
	"live2text/internal/services/burner"
	"live2text/internal/services/metrics"
	"live2text/internal/services/recognition"
	"live2text/internal/services/recognition/components"
	speechwrapper "live2text/internal/services/speech_wrapper"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args[1:]); err != nil {
		slog.Error("Application error", "error", err) //nolint:sloglint
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Initialize(os.Stderr, args)
	if err != nil {
		return fmt.Errorf("cannot initialize config: %w", err)
	}

	sl, l := newLoggers(cfg.LogLevel)

	s, sc, aw, tm, sm, err := newServices(ctx, sl, cfg)
	if err != nil {
		return fmt.Errorf("cannot create services: %w", err)
	}
	defer sc.Close()
	defer aw.Close()

	if err = checkPortAvailable(cfg.AppAddress); err != nil {
		return fmt.Errorf("app server port check failed: %w", err)
	}
	if err = checkPortAvailable(cfg.PprofAddress); err != nil {
		return fmt.Errorf("pprof server port check failed: %w", err)
	}

	server := newServer(ctx, cfg.AppAddress, api.NewHandler(sl, s), l)
	defer server.Close()
	go func() {
		sl.InfoContext(ctx, "Starting the API server", "address", cfg.AppAddress)
		serverErr := server.ListenAndServe()
		if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			sl.ErrorContext(ctx, "Failed to listen and serve", "error", serverErr)
			cancel()
		}
	}()

	pprofServer := newServer(ctx, cfg.PprofAddress, newPprofHandler(), l)
	go func() {
		sl.InfoContext(ctx, "Starting pprof server", "address", pprofServer.Addr)
		pprofServerErr := pprofServer.ListenAndServe()
		if pprofServerErr != nil && !errors.Is(pprofServerErr, http.ErrServerClosed) {
			sl.ErrorContext(ctx, "Pprof server error", "error", pprofServerErr)
			cancel()
		}
	}()

	// Wait for the interrupt signal
	<-ctx.Done()
	sl.InfoContext(ctx, "Received shutdown signal", "signal", ctx.Err())

	// Shutdown servers and clean up resource
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		sl.ErrorContext(ctx, "Failed to shut down the app server gracefully", "error", err)
	} else {
		sl.InfoContext(ctx, "Server shutdown completed successfully")
	}

	if err = pprofServer.Shutdown(shutdownCtx); err != nil {
		sl.ErrorContext(ctx, "Failed to shut down the pprof server gracefully", "error", err)
	} else {
		sl.InfoContext(ctx, "Pprof server shutdown completed successfully")
	}

	sl.InfoContext(ctx, "Closing socket manager")
	if err = sm.Close(); err != nil {
		sl.ErrorContext(ctx, "Failed to close socket manager", "error", err)
	}

	sl.InfoContext(ctx, "Waiting for background tasks to complete")
	tm.Wait()

	sl.InfoContext(ctx, "Application shutdown complete")

	return nil
}

func newServices(
	ctx context.Context,
	sl *slog.Logger,
	cfg *config.Config,
) (services.Services, speechwrapper.SpeechClient, audiowrapper.Audio, *background.TaskManager, *background.SocketManager, error) {
	sc, err := speechwrapper.NewSpeechClient(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("cannot create a speech client: %w", err)
	}

	aw, err := audiowrapper.NewAudio()
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("cannot create a new audio wrapper: %w", err)
	}

	tml := serviceLogger(sl, "task_manager")
	sml := serviceLogger(sl, "socket_manager")
	tm := background.NewTaskManager(ctx, tml)
	sm := background.NewSocketManager(sml)

	al := serviceLogger(sl, "audio")
	bl := serviceLogger(sl, "burner")
	rl := serviceLogger(sl, "recognition")
	bcl := componentLogger(rl, "burner")
	rcl := componentLogger(rl, "recognizer")
	scl := componentLogger(rl, "socket")
	ocl := componentLogger(rl, "output")
	m := metrics.NewMetrics(tm.TotalRunningTasks, sm.TotalOpenSockets)
	a := audio.NewAudio(al, m, aw)
	b := burner.NewBurner(bl, m)
	bc := components.NewBurnerComponent(bcl, b, cfg.OutputDir)
	rc := components.NewRecognizerComponent(rcl, m, sc)
	scc := components.NewSocketComponent(scl, sm)
	oc := components.NewOutputComponent(ocl, cfg.OutputDir, cfg.ConsoleWriter)
	tf := recognition.NewTaskFactory(rl, a, bc, rc, scc, oc)
	r := recognition.NewRecognition(rl, tm, tf)

	btl := serviceLogger(sl, "btt")
	bt := newBtt(btl, a, r, cfg)

	s := services.NewServices(a, aw, b, r, m, bt)

	return s, sc, aw, tm, sm, nil
}

func newLoggers(level slog.Level) (*slog.Logger, *log.Logger) {
	noColor := !term.IsTerminal(int(os.Stderr.Fd()))

	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:   level,
		NoColor: noColor,
	})

	return slog.New(handler), slog.NewLogLogger(handler, level)
}

func newBtt(logger *slog.Logger, audio audio.Audio, recognition recognition.Recognition, cfg *config.Config) btt.Btt {
	httpClient := httpclient.NewClient(logger, cfg.BttAddress, nil)
	client := bttclient.NewClient(httpClient, cfg.AppName)
	s := storage.NewStorage(httpClient)

	renderer := tmpl.NewRenderer(cfg.AppName, cfg.AppAddress, cfg.BttAddress, env.IsDebugMode())

	sc := btt.NewSettingsComponent(client, s)
	dc := btt.NewDeviceComponent(audio, client, renderer, sc)
	lc := btt.NewLanguageComponent(sc)
	vmc := btt.NewViewModeComponent(client, sc)
	fc := btt.NewFloatingComponent(logger, recognition, client, s, renderer, sc)
	lic := btt.NewListeningComponent(logger, recognition, client, s, renderer, dc, lc, vmc, fc)

	return btt.NewBtt(nil, lic, dc, lc, vmc, fc)
}

func checkPortAvailable(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("address %s is occupied: %w", address, err)
	}

	if err = l.Close(); err != nil {
		return fmt.Errorf("failed to close test listener: %w", err)
	}

	return nil
}

func newServer(ctx context.Context, address string, handler http.Handler, logger *log.Logger) *http.Server {
	if env.IsDebugMode() {
		return &http.Server{
			Addr:              address,
			Handler:           handler,
			ErrorLog:          logger,
			ReadTimeout:       60 * time.Minute,
			WriteTimeout:      60 * time.Minute,
			IdleTimeout:       60 * time.Minute,
			ReadHeaderTimeout: 60 * time.Minute,
			BaseContext:       func(_ net.Listener) context.Context { return ctx },
		}
	}

	return &http.Server{
		Addr:              address,
		Handler:           handler,
		ErrorLog:          logger,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}
}

func newPprofHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	return mux
}

func serviceLogger(base *slog.Logger, service string) *slog.Logger {
	return base.With("service", service)
}

func componentLogger(base *slog.Logger, component string) *slog.Logger {
	return base.With("component", component)
}
