package recognition_test

import (
	"log/slog"
	"testing"

	"github.com/gojuno/minimock/v3"

	"github.com/roboloop/live2text/internal/background"
	"github.com/roboloop/live2text/internal/services/audio"
	"github.com/roboloop/live2text/internal/services/recognition"
	"github.com/roboloop/live2text/internal/services/recognition/components"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func setupRecognition(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, a *audio.AudioMock, tf *recognition.TaskFactoryMock),
	l *slog.Logger,
) (recognition.Recognition, *background.TaskManager) {
	mc := minimock.NewController(t)
	a := audio.NewAudioMock(mc)
	tf := recognition.NewTaskFactoryMock(mc)

	if l == nil {
		l = logger.NilLogger
	}
	tm := background.NewTaskManager(t.Context(), l)

	if setupMocks != nil {
		setupMocks(mc, a, tf)
	}

	return recognition.NewRecognition(l, tm, tf), tm
}

func setupTaskFactory(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, a *audio.AudioMock, bc *components.BurnerComponentMock, rc *components.RecognizerComponentMock, scc *components.SocketComponentMock, oc *components.OutputComponentMock),
	l *slog.Logger,
) recognition.TaskFactory {
	mc := minimock.NewController(t)
	a := audio.NewAudioMock(mc)
	bc := components.NewBurnerComponentMock(mc)
	rc := components.NewRecognizerComponentMock(mc)
	scc := components.NewSocketComponentMock(mc)
	oc := components.NewOutputComponentMock(mc)

	if setupMocks != nil {
		setupMocks(mc, a, bc, rc, scc, oc)
	}

	if l == nil {
		l = logger.NilLogger
	}

	return recognition.NewTaskFactory(l, a, bc, rc, scc, oc)
}
