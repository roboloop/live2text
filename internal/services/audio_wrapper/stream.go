package audiowrapper

import "github.com/gordonklaus/portaudio"

type stream struct {
	*portaudio.Stream
	buffer *[]int16
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

func (s *stream) Read() ([]int16, error) {
	err := s.Stream.Read()
	if err != nil {
		return nil, err
	}

	// Create a copy of the buffer to return
	result := make([]int16, len(*s.buffer))
	copy(result, *s.buffer)

	return result, nil
}
