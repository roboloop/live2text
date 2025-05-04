package audio_wrapper

import (
	"github.com/gordonklaus/portaudio"
)

type MockAudio struct {
	OpenStreamStream *MockStream
	OpenStreamError  error

	DefaultHostApiHostApiInfo *portaudio.HostApiInfo
	DefaultHostApiError       error

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

func (m *MockAudio) DefaultHostApi() (*portaudio.HostApiInfo, error) {
	return m.DefaultHostApiHostApiInfo, m.DefaultHostApiError
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
