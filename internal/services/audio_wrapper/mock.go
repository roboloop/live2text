package audiowrapper

import (
	"github.com/gordonklaus/portaudio"
)

type MockAudio struct {
	OpenStreamStream *MockStream
	OpenStreamError  error

	DefaultHostAPIHostAPIInfo *portaudio.HostApiInfo
	DefaultHostAPIError       error

	CloseError error
}

type MockStream struct {
	StartError error
	ReadError  error
	CloseError error
	StopError  error
}

func (m *MockAudio) OpenStream(portaudio.StreamParameters, ...any) (Stream, error) {
	return m.OpenStreamStream, m.OpenStreamError
}

func (m *MockAudio) DefaultHostAPI() (*portaudio.HostApiInfo, error) {
	return m.DefaultHostAPIHostAPIInfo, m.DefaultHostAPIError
}

func (m *MockAudio) Close() error {
	return m.CloseError
}

func (m *MockStream) Close() error {
	return m.CloseError
}

func (m *MockStream) Start() error {
	return m.StartError
}

func (m *MockStream) Stop() error {
	return m.StopError
}

func (m *MockStream) Read() error {
	return m.ReadError
}
