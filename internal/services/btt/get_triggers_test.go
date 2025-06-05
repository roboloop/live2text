package btt_test

import (
	"log/slog"

	"live2text/internal/config"
	"live2text/internal/services/audio"
	btthttp "live2text/internal/services/btt/http"
	"live2text/internal/services/recognition"
)

func newMocks() (*audio.MockAudio, *recognition.MockRecognition, *btthttp.MockClient, *config.Config) { //nolint:unparam
	cfg := &config.Config{
		AppAddress: "127.0.0.1:1234",
		BttAddress: "127.0.0.1:5678",
		Languages:  []string{"en-US"},
		LogLevel:   slog.LevelInfo,
	}

	// TODO:
	mockAudio := &audio.MockAudio{
		ListDeviceInfo:             nil,
		ListError:                  nil,
		FindInputDeviceDeviceInfo:  nil,
		FindInputDeviceError:       nil,
		ListenDeviceDeviceListener: nil,
		ListenDeviceError:          nil,
	}

	mockRecognition := &recognition.MockRecognition{
		StartID:         "",
		StartSocketPath: "",
		StartError:      nil,
		StopError:       nil,
		SubsText:        "",
		SubsError:       nil,
	}

	mockHTTPClient := &btthttp.MockClient{
		SendResponse: nil,
		SendError:    nil,
	}

	return mockAudio, mockRecognition, mockHTTPClient, cfg
}
