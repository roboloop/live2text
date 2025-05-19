package btt_test

import (
	"live2text/internal/config"
	"live2text/internal/services/audio"
	bttexec "live2text/internal/services/btt/exec"
	btthttp "live2text/internal/services/btt/http"
	"live2text/internal/services/recognition"
	"log/slog"
)

func newMocks() (*audio.MockAudio, *recognition.MockRecognition, *btthttp.MockClient, *bttexec.MockClient, *config.Config) {
	//logger := utils.NilLogger
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
		StartId:         "",
		StartSocketPath: "",
		StartError:      nil,
		StopError:       nil,
		SubsText:        "",
		SubsError:       nil,
	}

	mockHttpClient := &btthttp.MockClient{
		SendResponse: nil,
		SendError:    nil,
	}

	mockExecClient := &bttexec.MockClient{
		ExecResponse: nil,
		ExecError:    nil,
	}

	//b := btt.NewBtt(logger, mockAudio, mockRecognition, mockHttpClient, mockExecClient, cfg)

	return mockAudio, mockRecognition, mockHttpClient, mockExecClient, cfg
}
