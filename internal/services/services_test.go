package services_test

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/audio"
	audiowrapper "github.com/roboloop/live2text/internal/services/audio_wrapper"
	"github.com/roboloop/live2text/internal/services/btt"
	"github.com/roboloop/live2text/internal/services/burner"
	"github.com/roboloop/live2text/internal/services/metrics"
	"github.com/roboloop/live2text/internal/services/recognition"
)

func TestServices(t *testing.T) {
	t.Parallel()

	t.Run("returns all services", func(t *testing.T) {
		t.Parallel()

		mc := minimock.NewController(t)
		a := audio.NewAudioMock(mc)
		aw := audiowrapper.NewAudioMock(mc)
		b := burner.NewBurnerMock(mc)
		r := recognition.NewRecognitionMock(mc)
		m := metrics.NewMetricsMock(mc)
		b2 := btt.NewBttMock(mc)

		s := services.NewServices(a, aw, b, r, m, b2)

		require.Equal(t, a, s.Audio())
		require.Equal(t, aw, s.AudioWrapper())
		require.Equal(t, b, s.Burner())
		require.Equal(t, r, s.Recognition())
		require.Equal(t, m, s.Metrics())
		require.Equal(t, b2, s.Btt())
	})
}
