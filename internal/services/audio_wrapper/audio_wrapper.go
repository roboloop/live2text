package audio_wrapper

import (
	"errors"
	"github.com/gordonklaus/portaudio"
)

type audio struct {
}

type stream struct {
	*portaudio.Stream
}

func NewAudio() (Audio, func() error, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, nil, errors.New("cannot initialize portaudio")
	}

	closeFn := func() error {
		return portaudio.Terminate()
	}

	return &audio{}, closeFn, nil
}

func (a *audio) OpenStream(p portaudio.StreamParameters, args ...any) (Stream, error) {
	s, err := portaudio.OpenStream(p, args...)
	return &stream{s}, err
}

func (a *audio) DefaultHostApi() (*portaudio.HostApiInfo, error) {
	return portaudio.DefaultHostApi()
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
