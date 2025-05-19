package audio

import (
	"context"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"time"
)

// ListenDevice listens the input device
// To stop listening cancel the context
func (a *audio) ListenDevice(ctx context.Context, deviceName string) (*DeviceListener, error) {
	var (
		device *portaudio.DeviceInfo
		err    error
		//defaultChunkSizeMs = 100
	)
	if device, err = a.FindInputDevice(deviceName); err != nil {
		return nil, fmt.Errorf("cannot find input device: %w", err)
	}

	return NewDeviceListener(a.logger, a.metrics, a.externalAudio, device), nil

	//ch := make(chan []int16, 1024)
	//errCh := make(chan error, 1)
	//info := &ListenerInfo{
	//	Device:      device,
	//	Channels:    1,
	//	SampleRate:  int(device.DefaultSampleRate),
	//	ChunkSizeMs: defaultChunkSizeMs,
	//	Ch:          ch,
	//	ErrCh:       errCh,
	//}
	//
	//go func() {
	//	defer close(ch)
	//	defer close(errCh)
	//	if err = a.listenDevice(ctx, info, ch); err != nil {
	//		errCh <- err
	//	}
	//	a.logger.InfoContext(ctx, "[Audio] Shutting down...")
	//}()
	//
	//return info, nil
}

func (a *audio) listenDevice(ctx context.Context, info *ListenerInfo, ch chan []int16) error {
	bufferSize := int(info.Device.DefaultSampleRate) * info.ChunkSizeMs / int(time.Second.Milliseconds())
	buffer := make([]int16, bufferSize*info.Channels)

	stream, err := a.externalAudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   info.Device,
			Channels: info.Channels,
			Latency:  info.Device.DefaultLowInputLatency,
		},
		SampleRate:      float64(info.SampleRate),
		FramesPerBuffer: bufferSize,
	}, buffer)
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

		select {
		case <-ctx.Done():
			a.logger.InfoContext(ctx, "shutdown of audio reader")
			return nil
		case ch <- clone:
		default:
			a.logger.ErrorContext(ctx, "The channel is full, segment was dropped")
		}
	}
}
