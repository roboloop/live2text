package audio

import (
	"fmt"
)

func (a *audio) ListOfNames() ([]string, error) {
	deviceInfos, err := a.externalAudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("cannot get a list of devices: %w", err)
	}

	var devices []string
	for _, device := range deviceInfos {
		if device.MaxInputChannels >= 1 {
			devices = append(devices, device.Name)
		}
	}

	return devices, nil
}
