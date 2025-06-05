package metrics_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"live2text/internal/services/metrics"
)

func TestMetrics(t *testing.T) {
	m := metrics.NewMetrics()
	m.AddBytesSentToGoogleSpeech(10)
	m.AddBytesWrittenOnDisk(20)
	m.AddBytesReadFromAudio(30)

	buf := bytes.NewBuffer([]byte{})
	m.WritePrometheus(buf)
	s := buf.String()

	if !strings.Contains(s, fmt.Sprintf("%s %d", "recognizer_bytes_sent_to_google_speech", 10)) {
		t.Errorf("Prometheus format does not contain %s", "recognizer_bytes_sent_to_google_speech")
	}
	if !strings.Contains(s, fmt.Sprintf("%s %d", "recognizer_bytes_written_on_disk", 20)) {
		t.Errorf("Prometheus format does not contain %s", "recognizer_bytes_written_on_disk")
	}
	if !strings.Contains(s, fmt.Sprintf("%s %d", "recognizer_bytes_read_from_audio", 30)) {
		t.Errorf("Prometheus format does not contain %s", "recognizer_bytes_read_from_audio")
	}
}
