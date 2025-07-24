package metrics

import (
	"io"

	externalmetrics "github.com/VictoriaMetrics/metrics"
)

const (
	MetricBytesSentToGoogleSpeech        = "app_bytes_sent_to_google_speech"
	MetricBytesWrittenOnDisk             = "app_bytes_written_on_disk"
	MetricBytesReadFromAudio             = "app_bytes_read_from_audio"
	MetricMillisecondsSentToGoogleSpeech = "app_milliseconds_sent_to_google_speech"
	MetricConnectsToGoogleSpeech         = "app_connections_to_google_speech"
	MetricTotalRunningTasks              = "app_total_running_tasks"
	MetricTotalOpenSockets               = "app_total_open_sockets"
)

type metrics struct {
	set *externalmetrics.Set

	bytesSentToGoogleSpeechCounter        *externalmetrics.Counter
	bytesWrittenOnDiskCounter             *externalmetrics.Counter
	bytesReadFromAudioCounter             *externalmetrics.Counter
	millisecondsSentToGoogleSpeechCounter *externalmetrics.Counter
	connectsToGoogleSpeechCounter         *externalmetrics.Counter
	totalRunningTasksGauge                *externalmetrics.Gauge
	totalOpenSocketsGauge                 *externalmetrics.Gauge
}

func NewMetrics(totalRunningTasks, totalOpenSockets func() int) Metrics {
	totalRunningTasksFn := func() float64 {
		return float64(totalRunningTasks())
	}
	totalOpenSocketsFn := func() float64 {
		return float64(totalOpenSockets())
	}

	set := externalmetrics.NewSet()
	var (
		bytesSentToGoogleSpeechCounter        = set.NewCounter(MetricBytesSentToGoogleSpeech)
		bytesWrittenOnDiskCounter             = set.NewCounter(MetricBytesWrittenOnDisk)
		bytesReadFromAudioCounter             = set.NewCounter(MetricBytesReadFromAudio)
		millisecondsSentToGoogleSpeechCounter = set.NewCounter(MetricMillisecondsSentToGoogleSpeech)
		connectsToGoogleSpeechCounter         = set.NewCounter(MetricConnectsToGoogleSpeech)
		totalRunningTasksGauge                = set.NewGauge(MetricTotalRunningTasks, totalRunningTasksFn)
		totalOpenSocketsGauge                 = set.NewGauge(MetricTotalOpenSockets, totalOpenSocketsFn)
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
	m.bytesReadFromAudioCounter.Add(bytes)
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
