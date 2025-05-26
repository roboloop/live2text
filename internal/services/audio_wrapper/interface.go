package audiowrapper

import "github.com/gordonklaus/portaudio"

type Audio interface {
	OpenStream(portaudio.StreamParameters, ...any) (Stream, error)
	DefaultHostAPI() (*portaudio.HostApiInfo, error)
	Close() error
}

type Stream interface {
	Close() error
	Start() error
	Stop() error
	Read() error
}
