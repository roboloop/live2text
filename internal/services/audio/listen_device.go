package audio

import (
	"context"
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// ListenDevice returns device listener.
func (a *audio) ListenDevice(_ context.Context, deviceName string) (*DeviceListener, error) {
	var (
		device *portaudio.DeviceInfo
		err    error
	)
	if device, err = a.FindInputDevice(deviceName); err != nil {
		return nil, fmt.Errorf("cannot find input device: %w", err)
	}

	return NewDeviceListener(a.logger, a.metrics, a.externalAudio, device), nil
}
