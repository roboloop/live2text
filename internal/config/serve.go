package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/roboloop/live2text/internal/env"
)

type Serve struct {
	AppAddress   string
	PprofAddress string
	BttAddress   string

	OutputDir     string
	LogLevel      slog.Level
	ConsoleWriter io.Writer
}

func (s *Serve) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("app-address", s.AppAddress),
		slog.String("pprof-address", s.PprofAddress),
		slog.String("btt-address", s.BttAddress),
		slog.String("output-dir", s.OutputDir),
		slog.String("log-level", s.LogLevel.String()),
		slog.String("on-console", strconv.FormatBool(s.ConsoleWriter == os.Stdout)),
	)
}

func ParseServe(output io.Writer, args []string) (*Serve, error) {
	var (
		appHost   string
		appPort   string
		pprofHost string
		pprofPort string
		bttHost   string
		bttPort   string

		outputDir string
		logLevel  string
		onConsole bool
	)

	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(output)
	fs.StringVar(&appHost, "app-host", "127.0.0.1", "App server host")
	fs.StringVar(&appPort, "app-port", "8080", "App server port")
	fs.StringVar(&pprofHost, "pprof-host", "127.0.0.1", "Pprof server host")
	fs.StringVar(&pprofPort, "pprof-port", "8081", "Pprof server port")
	fs.StringVar(&bttHost, "btt-host", "127.0.0.1", "BetterTouchTool webserver host")
	fs.StringVar(&bttPort, "btt-port", "44444", "BetterTouchTool webserver port")
	fs.StringVar(
		&outputDir,
		"output-dir",
		filepath.Join(os.TempDir(), env.AppName),
		"Dir to the produced audio and text",
	)
	fs.StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn, error")
	fs.BoolVar(&onConsole, "on-console", false, "Print subtitles on the console")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, ErrHelp
		}

		return nil, fmt.Errorf("cannot parse arguments: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return nil, fmt.Errorf("cannot create the output dir: %w", err)
	}

	return &Serve{
		AppAddress:    net.JoinHostPort(appHost, appPort),
		PprofAddress:  net.JoinHostPort(pprofHost, pprofPort),
		BttAddress:    net.JoinHostPort(bttHost, bttPort),
		OutputDir:     outputDir,
		LogLevel:      parseLogLevel(logLevel),
		ConsoleWriter: consoleWriter(onConsole),
	}, nil
}

func consoleWriter(onConsole bool) io.Writer {
	if onConsole {
		return os.Stdout
	}
	return io.Discard
}
