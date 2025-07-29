package burner

import (
	"context"
	"fmt"
	"io"

	"github.com/youpy/go-wav"
)

func (b *burner) Burn(ctx context.Context, w io.Writer, input <-chan []int16, channels, sampleRate int) error {
	const bitsPerSample = uint16(16)
	const bitsPerByte = 8
	var samples []wav.Sample

	for {
		select {
		case <-ctx.Done():
			b.logger.InfoContext(ctx, "Writing samples...", "total", len(samples))

			//nolint:gosec // the values always small, so just ignore type conversion
			writer := wav.NewWriter(w, uint32(len(samples)), uint16(channels), uint32(sampleRate), bitsPerSample)
			if err := writer.WriteSamples(samples); err != nil {
				return fmt.Errorf("cannot write samples: %w", err)
			}

			b.logger.InfoContext(ctx, "Writing samples is done!")
			b.metrics.AddBytesWrittenOnDisk(len(samples) * int(bitsPerSample/bitsPerByte))
			return nil
		case value := <-input:
			if channels == 1 {
				samples = append(samples, int16ToSample(value)...)
			} else {
				samples = append(samples, int16ToSampleInStereo(value)...)
			}
		}
	}
}

func int16ToSample(buffer []int16) []wav.Sample {
	samples := make([]wav.Sample, len(buffer))
	for i, v := range buffer {
		samples[i] = wav.Sample{Values: [2]int{int(v), 0}}
	}
	return samples
}

func int16ToSampleInStereo(buffer []int16) []wav.Sample {
	samples := make([]wav.Sample, len(buffer)/2)
	for i := 0; i < len(buffer); i += 2 {
		samples[i/2] = wav.Sample{Values: [2]int{int(buffer[i]), int(buffer[i+1])}}
	}
	return samples
}
