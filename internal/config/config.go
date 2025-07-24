package config

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
)

type Config struct {
	AppAddress   string
	PprofAddress string
	BttAddress   string

	AppName       string
	OutputDir     string
	LogLevel      slog.Level
	ConsoleWriter io.Writer
}

type BttConfig struct {
	AppAddress string
	BttAddress string

	AppName   string
	Languages []string
	LogLevel  slog.Level
	Clear     bool
}

func Initialize(output io.Writer, args []string) (*Config, error) {
	var (
		appHost   string
		appPort   string
		pprofHost string
		pprofPort string
		bttHost   string
		bttPort   string

		appName   string
		outputDir string
		logLevel  string
		onConsole bool
	)

	// TODO: No flag duplication
	fs := flag.NewFlagSet("live2text", flag.ContinueOnError)
	fs.SetOutput(output)
	fs.StringVar(&appHost, "app-host", "127.0.0.1", "App server host")
	fs.StringVar(&appPort, "app-port", "8080", "App server port")
	fs.StringVar(&pprofHost, "pprof-host", "127.0.0.1", "Pprof server host")
	fs.StringVar(&pprofPort, "pprof-port", "8081", "Pprof server port")
	fs.StringVar(&bttHost, "btt-host", "127.0.0.1", "BetterTouchTool webserver host")
	fs.StringVar(&bttPort, "btt-port", "44444", "BetterTouchTool webserver port")
	fs.StringVar(&appName, "app-name", "Live2Text", "App name in specified in labels")
	fs.StringVar(&outputDir, "output-dir", os.TempDir()+"live2text", "Dir to the produced audio and text")
	fs.StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn, error")
	fs.BoolVar(&onConsole, "on-console", false, "Print subtitles on the console")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("cannot parse arguments: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return nil, fmt.Errorf("cannot create the output dir: %w", err)
	}

	return &Config{
		AppAddress:    net.JoinHostPort(appHost, appPort),
		PprofAddress:  net.JoinHostPort(pprofHost, pprofPort),
		BttAddress:    net.JoinHostPort(bttHost, bttPort),
		AppName:       appName,
		OutputDir:     outputDir,
		LogLevel:      parseLogLevel(logLevel),
		ConsoleWriter: consoleWriter(onConsole),
	}, nil
}

func InitializeBtt(output io.Writer, args []string) (*BttConfig, error) {
	var (
		appHost string
		appPort string
		bttHost string
		bttPort string

		appName   string
		languages string
		logLevel  string
		clearAll  bool
	)

	fs := flag.NewFlagSet("btt", flag.ContinueOnError)
	fs.SetOutput(output)
	fs.StringVar(&appHost, "app-host", "127.0.0.1", "App server host")
	fs.StringVar(&appPort, "app-port", "8080", "App server port")
	fs.StringVar(&bttHost, "btt-host", "127.0.0.1", "BetterTouchTool webserver host")
	fs.StringVar(&bttPort, "btt-port", "44444", "BetterTouchTool webserver port")
	fs.StringVar(&appName, "app-name", "Live2Text", "App name in specified in labels")
	fs.StringVar(&languages, "languages", "en-US,es-ES,fr-FR,pt-BR,ru-RU,ja-JP,de-DE", "List of available languages")
	fs.StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn, error")
	fs.BoolVar(&clearAll, "clear", false, "Clear all btt triggers")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("cannot parse arguments: %w", err)
	}

	return &BttConfig{
		AppAddress: net.JoinHostPort(appHost, appPort),
		BttAddress: net.JoinHostPort(bttHost, bttPort),
		AppName:    appName,
		Languages:  strings.Split(languages, ","),
		LogLevel:   parseLogLevel(logLevel),
		Clear:      clearAll,
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

func consoleWriter(onConsole bool) io.Writer {
	if onConsole {
		return os.Stdout
	}
	return io.Discard
}
