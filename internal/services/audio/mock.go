package audio

import (
	"context"

	"github.com/gordonklaus/portaudio"
)

type MockAudio struct {
	ListDeviceInfo []*portaudio.DeviceInfo
	ListError      error

	FindInputDeviceDeviceInfo *portaudio.DeviceInfo
	FindInputDeviceError      error

	ListenDeviceDeviceListener *DeviceListener
	ListenDeviceError          error
}

func (m *MockAudio) List() ([]*portaudio.DeviceInfo, error) {
	return m.ListDeviceInfo, m.ListError
}

func (m *MockAudio) FindInputDevice(string) (*portaudio.DeviceInfo, error) {
	return m.FindInputDeviceDeviceInfo, m.FindInputDeviceError
}

func (m *MockAudio) ListenDevice(context.Context, string) (*DeviceListener, error) {
	return m.ListenDeviceDeviceListener, m.ListenDeviceError
}
