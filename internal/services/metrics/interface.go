package metrics

import (
	"io"
)

//go:generate minimock -g -i Metrics -s _mock.go -o .

type Metrics interface {
	AddBytesSentToGoogleSpeech(bytes int)
	AddBytesWrittenOnDisk(bytes int)
	AddBytesReadFromAudio(bytes int)
	AddMillisecondsSentToGoogleSpeech(ms int)
	AddConnectionsToGoogleSpeech(n int)

	WritePrometheus(w io.Writer)
}
