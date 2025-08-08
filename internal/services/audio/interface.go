package audio

import (
	"context"

	audiowrapper "github.com/roboloop/live2text/internal/services/audio_wrapper"
)

//go:generate minimock -g -i Audio -s _mock.go -o .
//go:generate minimock -g -i DeviceListener -s _mock.go -o .

type Audio interface {
	ListOfNames() ([]string, error)
	ListenDevice(deviceName string) (DeviceListener, error)
	FindInputDevice(deviceName string) (*audiowrapper.DeviceInfo, error)
}

// DeviceListener defines the contract for a device that can be listened to.
type DeviceListener interface {
	Listen(ctx context.Context) error
	GetChannel() <-chan []int16
	GetParameters() *audiowrapper.StreamParameters
}
