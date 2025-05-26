package audio

import (
	"errors"
	"fmt"

	"github.com/gordonklaus/portaudio"
)

func (a *audio) FindInputDevice(deviceName string) (*portaudio.DeviceInfo, error) {
	hostAPI, err := a.externalAudio.DefaultHostAPI()
	if err != nil {
		return nil, fmt.Errorf("cannot list host apis: %w", err)
	}

	var device *portaudio.DeviceInfo
	for _, d := range hostAPI.Devices {
		if d.Name == deviceName {
			device = d
		}
	}
	if device == nil {
		return nil, errors.New("device not found")
	}
	if device.MaxInputChannels <= 0 {
		return nil, errors.New("device hasn't input channels")
	}

	return device, nil
}
