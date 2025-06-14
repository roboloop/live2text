package metrics

import (
	"io"

	externalmetrics "github.com/VictoriaMetrics/metrics"
)

const (
	MetricBytesSentToGoogleSpeech        = "recognizer_bytes_sent_to_google_speech"
	MetricBytesWrittenOnDisk             = "recognizer_bytes_written_on_disk"
	MetricBytesReadFromAudio             = "recognizer_bytes_read_from_audio"
	MetricMillisecondsSentToGoogleSpeech = "recognizer_milliseconds_sent_to_google_speech"
	MetricConnectsToGoogleSpeech         = "recognizer_connections_to_google_speech"
	MetricTotalRunningTasks              = "recognizer_total_running_tasks"
	MetricTotalOpenSockets               = "recognizer_total_open_sockets"
)

type metrics struct {
	set *externalmetrics.Set

	bytesSentToGoogleSpeechCounter        *externalmetrics.Counter
	bytesWrittenOnDiskCounter             *externalmetrics.Counter
	bytesReadFromAudio                    *externalmetrics.Counter
	millisecondsSentToGoogleSpeechCounter *externalmetrics.Counter
	connectsToGoogleSpeechCounter         *externalmetrics.Counter
	totalRunningTasksGauge                *externalmetrics.Gauge
	totalOpenSocketsGauge                 *externalmetrics.Gauge
}

func NewMetrics(totalRunningTasks, totalOpenSockets func() float64) Metrics {
	set := externalmetrics.NewSet()
	var (
		bytesSentToGoogleSpeechCounter        = set.NewCounter(MetricBytesSentToGoogleSpeech)
		bytesWrittenOnDiskCounter             = set.NewCounter(MetricBytesWrittenOnDisk)
		bytesReadFromAudioCounter             = set.NewCounter(MetricBytesReadFromAudio)
		millisecondsSentToGoogleSpeechCounter = set.NewCounter(MetricMillisecondsSentToGoogleSpeech)
		connectsToGoogleSpeechCounter         = set.NewCounter(MetricConnectsToGoogleSpeech)
		totalRunningTasksGauge                = set.NewGauge(MetricTotalRunningTasks, totalRunningTasks)
		totalOpenSocketsGauge                 = set.NewGauge(MetricTotalOpenSockets, totalOpenSockets)
	)

	return &metrics{
		set,

		bytesSentToGoogleSpeechCounter,
		bytesWrittenOnDiskCounter,
		bytesReadFromAudioCounter,
		millisecondsSentToGoogleSpeechCounter,
		connectsToGoogleSpeechCounter,
		totalRunningTasksGauge,
		totalOpenSocketsGauge,
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

func (m *metrics) AddMillisecondsSentToGoogleSpeech(milliseconds int) {
	m.millisecondsSentToGoogleSpeechCounter.Add(milliseconds)
}

func (m *metrics) AddConnectionsToGoogleSpeech(n int) {
	m.connectsToGoogleSpeechCounter.Add(n)
}

func (m *metrics) WritePrometheus(w io.Writer) {
	externalmetrics.WriteProcessMetrics(w)
	externalmetrics.WriteFDMetrics(w)
	m.set.WritePrometheus(w)
}
