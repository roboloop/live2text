package burner

import (
	"context"
	"io"
)

//go:generate minimock -g -i Burner -s _mock.go -o .

type Burner interface {
	Burn(ctx context.Context, w io.Writer, input <-chan []int16, channels, sampleRate int) error
}
