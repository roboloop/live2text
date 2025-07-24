package audio

import (
	"context"
	"fmt"
	"log/slog"

	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/services/metrics"
)

type deviceListener struct {
	logger  *slog.Logger
	metrics metrics.Metrics

	externalAudio audiowrapper.Audio
	deviceInfo    *audiowrapper.DeviceInfo
	parameters    *audiowrapper.StreamParameters
	ch            chan []int16
}

func NewDeviceListener(
	logger *slog.Logger,
	metrics metrics.Metrics,
	externalAudio audiowrapper.Audio,
	deviceInfo *audiowrapper.DeviceInfo,
) DeviceListener {
	const (
		defaultChannels    = 1
		defaultChunkSizeMs = 100
		channelBufferSize  = 1024
	)
	// We have to fit in 10mb per connection (5 min default).
	// The size per minute is: 48000 * 2 * 60 / 1024 = 5.625mb (if bitrate is 48000)
	sampleRate := int(deviceInfo.DefaultSampleRate / 2)

	return &deviceListener{
		logger:        logger.With("component", "device_listener", "device", deviceInfo.Name),
		metrics:       metrics,
		externalAudio: externalAudio,
		deviceInfo:    deviceInfo,
		parameters: &audiowrapper.StreamParameters{
			Channels:    defaultChannels,
			ChunkSizeMs: defaultChunkSizeMs,
			SampleRate:  sampleRate,
		},
		ch: make(chan []int16, channelBufferSize),
	}
}

func (dl *deviceListener) GetChannel() <-chan []int16 {
	return dl.ch
}

func (dl *deviceListener) GetParameters() *audiowrapper.StreamParameters {
	return dl.parameters
}

func (dl *deviceListener) Listen(ctx context.Context) error {
	stream, err := dl.externalAudio.StreamDevice(dl.deviceInfo, &audiowrapper.StreamParameters{
		Channels:    dl.parameters.Channels,
		SampleRate:  dl.parameters.SampleRate,
		ChunkSizeMs: dl.parameters.ChunkSizeMs,
	})
	if err != nil {
		return fmt.Errorf("cannot open the stream: %w", err)
	}
	defer func() {
		if errStream := stream.Close(); errStream != nil {
			dl.logger.ErrorContext(ctx, "Cannot close stream", "error", errStream)
		}
	}()

	if err = stream.Start(); err != nil {
		return fmt.Errorf("cannot start the stream: %w", err)
	}
	defer func() {
		if errStream := stream.Stop(); errStream != nil {
			dl.logger.ErrorContext(ctx, "Cannot stop the stream", "error", errStream)
		}
	}()

	// TODO: incorrect finalizing
	for {
		data, errRead := stream.Read()
		if errRead != nil {
			return fmt.Errorf("cannot read the stream: %w", errRead)
		}

		dl.metrics.AddBytesReadFromAudio(len(data) * 2) // 2 bytes in int16

		select {
		case <-ctx.Done():
			dl.logger.InfoContext(ctx, "shutdown of audio reader")
			close(dl.ch)
			return nil
		case dl.ch <- data:

		default:
			dl.logger.ErrorContext(ctx, "The channel is full, the segment was dropped")
		}
	}
}
