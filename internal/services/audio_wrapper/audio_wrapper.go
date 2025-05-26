package audiowrapper

import (
	"errors"
	"sync"

	"github.com/gordonklaus/portaudio"
)

type audio struct {
	closeOnce sync.Once
	errClose  error
}

type stream struct {
	*portaudio.Stream
}

func NewAudio() (Audio, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, errors.New("cannot initialize portaudio")
	}

	return &audio{}, nil
}

func (a *audio) OpenStream(p portaudio.StreamParameters, args ...any) (Stream, error) {
	s, err := portaudio.OpenStream(p, args...)
	return &stream{s}, err
}

func (a *audio) DefaultHostAPI() (*portaudio.HostApiInfo, error) {
	return portaudio.DefaultHostApi()
}

func (a *audio) Close() error {
	a.closeOnce.Do(func() {
		a.errClose = portaudio.Terminate()
	})

	return a.errClose
}

func (s *stream) Close() error {
	return s.Stream.Close()
}

func (s *stream) Start() error {
	return s.Stream.Start()
}

func (s *stream) Stop() error {
	return s.Stream.Stop()
}

func (s *stream) Read() error {
	return s.Stream.Read()
}
