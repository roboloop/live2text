package audio

import (
	"fmt"
)

// ListenDevice returns a device listener that implements the DeviceListener interface.
func (a *audio) ListenDevice(deviceName string) (DeviceListener, error) {
	device, err := a.FindInputDevice(deviceName)
	if err != nil {
		return nil, fmt.Errorf("cannot find the input device: %w", err)
	}

	return NewDeviceListener(a.logger, a.metrics, a.externalAudio, device), nil
}
