package audio

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
)

func (a *audio) FindInputDevice(deviceName string) (*portaudio.DeviceInfo, error) {
	hostApi, err := a.externalAudio.DefaultHostApi()
	if err != nil {
		return nil, fmt.Errorf("cannot list host apis: %w", err)
	}

	var device *portaudio.DeviceInfo
	for _, d := range hostApi.Devices {
		if d.Name == deviceName {
			device = d
		}
	}
	if device == nil {
		return nil, fmt.Errorf("device not found")
	}
	if device.MaxInputChannels <= 0 {
		return nil, fmt.Errorf("device hasn't input channels")
	}

	return device, nil
}
