package audiowrapper

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

type DeviceInfo struct {
	Name              string
	DefaultSampleRate float64
	MaxInputChannels  int

	device *portaudio.DeviceInfo
}

type audio struct {
	closeOnce sync.Once
	errClose  error
}

func NewAudio() (Audio, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, errors.New("cannot initialize portaudio")
	}

	return &audio{}, nil
}

func (a *audio) Close() error {
	a.closeOnce.Do(func() {
		a.errClose = portaudio.Terminate()
	})

	return a.errClose
}

func (a *audio) Devices() ([]*DeviceInfo, error) {
	hostAPI, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, fmt.Errorf("cannot get default host api: %w", err)
	}

	var deviceInfos []*DeviceInfo
	for _, device := range hostAPI.Devices {
		deviceInfos = append(deviceInfos, &DeviceInfo{
			Name:              device.Name,
			DefaultSampleRate: device.DefaultSampleRate,
			MaxInputChannels:  device.MaxInputChannels,
			device:            device,
		})
	}

	return deviceInfos, nil
}

func (a *audio) StreamDevice(info *DeviceInfo, params *StreamParameters) (Stream, error) {
	if info.device == nil {
		return nil, errors.New("no portaudio device")
	}
	if params.SampleRate > int(info.device.DefaultSampleRate) {
		return nil, fmt.Errorf("cannot use sample rate greater than %f", info.device.DefaultSampleRate)
	}

	bufferSize := params.SampleRate * params.ChunkSizeMs / int(time.Second.Milliseconds())
	buffer := make([]int16, bufferSize*params.Channels)

	s, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   info.device,
			Channels: params.Channels,
			Latency:  info.device.DefaultHighInputLatency,
		},
		SampleRate:      float64(params.SampleRate),
		FramesPerBuffer: bufferSize,
	}, &buffer)
	if err != nil {
		return nil, fmt.Errorf("cannot open the stream: %w", err)
	}

	return &stream{Stream: s, buffer: &buffer}, nil
}
