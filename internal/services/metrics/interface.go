package metrics

import (
	"io"
)

type Metrics interface {
	AddBytesSentToGoogleSpeech(bytes int)
	AddBytesWrittenOnDisk(bytes int)
	AddBytesReadFromAudio(bytes int)
	AddMillisecondsSentToGoogleSpeech(ms int)
	AddConnectionsToGoogleSpeech(n int)

	WritePrometheus(w io.Writer)
}
