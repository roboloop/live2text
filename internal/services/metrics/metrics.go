package metrics

import (
	"io"

	externalmetrics "github.com/VictoriaMetrics/metrics"
)

var metricBytesSentToGoogleSpeech = "recognizer_bytes_sent_to_google_speech"
var metricBytesWrittenOnDisk = "recognizer_bytes_written_on_disk"
var metricBytesReadFromAudio = "recognizer_bytes_read_from_audio"

type metrics struct {
	set *externalmetrics.Set

	bytesSentToGoogleSpeechCounter *externalmetrics.Counter
	bytesWrittenOnDiskCounter      *externalmetrics.Counter
	bytesReadFromAudio             *externalmetrics.Counter
}

func NewMetrics() Metrics {
	set := externalmetrics.NewSet()
	bytesSentToGoogleSpeechCounter := set.NewCounter(metricBytesSentToGoogleSpeech)
	bytesWrittenOnDiskCounter := set.NewCounter(metricBytesWrittenOnDisk)
	bytesReadFromAudioCounter := set.NewCounter(metricBytesReadFromAudio)

	return &metrics{
		set,

		bytesSentToGoogleSpeechCounter,
		bytesWrittenOnDiskCounter,
		bytesReadFromAudioCounter,
	}
}

func (m *metrics) AddBytesSentToGoogleSpeech(bytes int) {
	m.bytesSentToGoogleSpeechCounter.Add(bytes)
}

func (m *metrics) AddBytesWrittenOnDisk(bytes int) {
	m.bytesWrittenOnDiskCounter.Add(bytes)
}

func (m *metrics) AddBytesReadFromAudio(bytes int) {
	m.bytesReadFromAudio.Add(bytes)
}

func (m *metrics) WritePrometheus(w io.Writer) {
	externalmetrics.WriteProcessMetrics(w)
	externalmetrics.WriteFDMetrics(w)
	m.set.WritePrometheus(w)
}
