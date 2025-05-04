package audio

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
)

func (a *audio) List() ([]*portaudio.DeviceInfo, error) {
	hostApi, err := a.externalAudio.DefaultHostApi()
	if err != nil {
		return nil, fmt.Errorf("cannot list host apis: %w", err)
	}

	var devices []*portaudio.DeviceInfo
	for _, device := range hostApi.Devices {
		if device.MaxInputChannels >= 1 {
			devices = append(devices, device)
		}
	}

	return devices, nil
}
