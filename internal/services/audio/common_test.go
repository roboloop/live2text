package audio_test

import (
	"log/slog"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services/audio"
	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/services/metrics"
	"live2text/internal/utils/logger"
)

func setupAudio(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock),
	l *slog.Logger,
) audio.Audio {
	mc := minimock.NewController(t)
	m := metrics.NewMetricsMock(mc)
	ea := audiowrapper.NewAudioMock(t)

	if setupMocks != nil {
		setupMocks(mc, m, ea)
	}
	if l == nil {
		l = logger.NilLogger
	}

	return audio.NewAudio(l, m, ea)
}

func setupDeviceListener(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock),
	l *slog.Logger,
	di *audiowrapper.DeviceInfo,
) audio.DeviceListener {
	mc := minimock.NewController(t)
	m := metrics.NewMetricsMock(mc)
	ea := audiowrapper.NewAudioMock(mc)

	if setupMocks != nil {
		setupMocks(mc, m, ea)
	}
	if l == nil {
		l = logger.NilLogger
	}
	if di == nil {
		di = &audiowrapper.DeviceInfo{}
	}

	return audio.NewDeviceListener(l, m, ea, di)
}
