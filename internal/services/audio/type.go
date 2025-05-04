package audio

import "github.com/gordonklaus/portaudio"

type ListenerInfo struct {
	Device      *portaudio.DeviceInfo
	Channels    int
	SampleRate  int
	ChunkSizeMs int

	Ch    <-chan []int16
	ErrCh <-chan error

	notImplemented
}

type notImplemented interface {
}
