package audio

import (
	"context"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"live2text/internal/services/audio_wrapper"
	"live2text/internal/services/metrics"
	"log/slog"
	"time"
)

type DeviceListener struct {
	logger        *slog.Logger
	metrics       metrics.Metrics
	externalAudio audio_wrapper.Audio

	deviceInfo *portaudio.DeviceInfo

	Channels    int
	SampleRate  int
	ChunkSizeMs int
	Ch          <-chan []int16
	ch          chan []int16
}

func NewDeviceListener(logger *slog.Logger, metrics metrics.Metrics, externalAudio audio_wrapper.Audio, deviceInfo *portaudio.DeviceInfo) *DeviceListener {
	const defaultChunkSizeMs = 100
	ch := make(chan []int16, 1024)
	return &DeviceListener{
		logger:        logger,
		metrics:       metrics,
		externalAudio: externalAudio,
		deviceInfo:    deviceInfo,

		Channels: 1,
		// TODO: we have to fit in 10mb per connection (5 min default). The size per minute is: 48000 * 2 * 60 / 1024 = 5.625mb (if bitrate is 48000)
		SampleRate:  int(deviceInfo.DefaultSampleRate / 2), // 48000 / 2.
		ChunkSizeMs: defaultChunkSizeMs,
		Ch:          ch,
		ch:          ch,
	}
}

func (dl *DeviceListener) Listen(ctx context.Context) error {
	bufferSize := dl.SampleRate * dl.ChunkSizeMs / int(time.Second.Milliseconds())
	buffer := make([]int16, bufferSize*dl.Channels)

	stream, err := dl.externalAudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   dl.deviceInfo,
			Channels: dl.Channels,
			Latency:  dl.deviceInfo.DefaultHighInputLatency,
		},
		SampleRate:      float64(dl.SampleRate),
		FramesPerBuffer: bufferSize,
	}, &buffer)
	if err != nil {
		return fmt.Errorf("cannot open the stream: %w", err)
	}
	defer stream.Close()

	if err = stream.Start(); err != nil {
		return fmt.Errorf("cannot start the stream: %w", err)
	}
	defer stream.Stop()

	for {
		if err = stream.Read(); err != nil {
			return fmt.Errorf("cannot read stream: %w", err)
		}
		clone := make([]int16, len(buffer))
		copy(clone, buffer)

		dl.metrics.AddBytesReadFromAudio(len(buffer) * 2)

		select {
		case <-ctx.Done():
			dl.logger.InfoContext(ctx, "shutdown of audio reader")
			close(dl.ch)
			return nil
		case dl.ch <- clone:
		default:
			dl.logger.ErrorContext(ctx, "The channel is full, segment was dropped")
		}
	}
}
