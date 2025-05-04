package audio

import (
	"context"
	"github.com/gordonklaus/portaudio"
)

type Audio interface {
	List() ([]*portaudio.DeviceInfo, error)
	FindInputDevice(deviceName string) (*portaudio.DeviceInfo, error)
	ListenDevice(ctx context.Context, deviceName string) (*DeviceListener, error)
}
