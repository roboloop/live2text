package metrics

import "io"

type Metrics interface {
	AddBytesSentToGoogleSpeech(bytes int)
	AddBytesWrittenOnDisk(bytes int)
	AddBytesReadFromAudio(bytes int)

	WritePrometheus(w io.Writer)
}
