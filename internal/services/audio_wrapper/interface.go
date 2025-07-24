package audiowrapper

//go:generate minimock -g -i Audio -s _mock.go -o .
//go:generate minimock -g -i Stream -s _mock.go -o .
type StreamParameters struct {
	Channels    int
	ChunkSizeMs int
	SampleRate  int
}

type Audio interface {
	Devices() ([]*DeviceInfo, error)
	StreamDevice(info *DeviceInfo, params *StreamParameters) (Stream, error)
	Close() error
}

type Stream interface {
	Close() error
	Start() error
	Stop() error
	Read() ([]int16, error)
}
