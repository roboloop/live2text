package audio

import (
	"errors"
	"fmt"

	audiowrapper "live2text/internal/services/audio_wrapper"
)

func (a *audio) FindInputDevice(deviceName string) (*audiowrapper.DeviceInfo, error) {
	devices, err := a.externalAudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("cannot get a list of devices: %w", err)
	}

	var device *audiowrapper.DeviceInfo
	for _, d := range devices {
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
