package burner

import (
	"context"
	"io"
)

type Burner interface {
	Burn(ctx context.Context, w io.Writer, input <-chan []int16, channels int, sampleRate int) error
}
