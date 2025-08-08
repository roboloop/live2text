package metrics_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/roboloop/live2text/internal/services/metrics"
)

func TestMetrics(t *testing.T) {
	t.Parallel()

	totalRunningTasks := func() int {
		return 1
	}
	totalOpenSockets := func() int {
		return 2
	}

	m := metrics.NewMetrics(totalRunningTasks, totalOpenSockets)
	m.AddBytesSentToGoogleSpeech(10)
	m.AddBytesWrittenOnDisk(20)
	m.AddBytesReadFromAudio(30)
	m.AddMillisecondsSentToGoogleSpeech(40)
	m.AddConnectionsToGoogleSpeech(50)

	buf := bytes.NewBuffer([]byte{})
	m.WritePrometheus(buf)
	s := buf.String()

	assert.Contains(t, s, fmt.Sprintf("%s %d", "app_bytes_sent_to_google_speech", 10))
	assert.Contains(t, s, fmt.Sprintf("%s %d", "app_bytes_written_on_disk", 20))
	assert.Contains(t, s, fmt.Sprintf("%s %d", "app_bytes_read_from_audio", 30))
	assert.Contains(t, s, fmt.Sprintf("%s %d", "app_milliseconds_sent_to_google_speech", 40))
	assert.Contains(t, s, fmt.Sprintf("%s %d", "app_connections_to_google_speech", 50))

	assert.Contains(t, s, fmt.Sprintf("%s %d", "app_total_running_tasks", 1))
	assert.Contains(t, s, fmt.Sprintf("%s %d", "app_total_open_sockets", 2))
}
