package config

import (
	"flag"
	"fmt"
)

type Config2 struct {
	Devices      []string
	SocketOutput string
	Languages    []string
}

type Config struct {
	Host string
	Port string
}

func Initialize(args []string) (*Config, error) {
	var (
		host string
		port string
	)

	fs := flag.NewFlagSet("live2text", flag.ContinueOnError)
	fs.StringVar(&host, "host", "127.0.0.1", "Host address")
	fs.StringVar(&port, "port", "8000", "Port address")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("cannot parse agruments: %w", err)
	}

	return &Config{host, port}, nil
}

func Initialize2(args []string) (*Config2, error) {
	var (
		devices      []string
		socketOutput string
		languages    []string
	)

	fs := flag.NewFlagSet("live2text", flag.ContinueOnError)
	fs.Func("device", "Device name (multiple allowed)", func(deviceName string) error {
		devices = append(devices, deviceName)
		return nil
	})
	fs.StringVar(&socketOutput, "output", "/tmp/live2text", "socket path for the output")
	fs.Func("language", "Language of speech (multiple allowed)", func(code string) error {
		//if !validation.ValidateCode(code) {
		//	return fmt.Errorf("%s is not valid code", code)
		//}
		languages = append(languages, code)
		return nil
	})

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("cannot parse agruments: %w", err)
	}

	if len(devices) == 0 {
		devices = append(devices, "Loopback Audio")
	}

	if len(languages) == 0 {
		languages = append(languages, "en-US")
	}

	return &Config2{
		Devices:      devices,
		SocketOutput: socketOutput,
		Languages:    languages,
	}, nil
}
