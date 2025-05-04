package speech_wrapper

import (
	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"context"
	"fmt"
)

type speechClient struct {
	c *speech.Client
}

func NewClient(ctx context.Context) (Client, error) {
	sc, err := speech.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create speech client: %w", err)
	}

	return &speechClient{sc}, nil
}

func (sc *speechClient) StreamingRecognize(ctx context.Context) (speechpb.Speech_StreamingRecognizeClient, error) {
	stream, err := sc.c.StreamingRecognize(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not streaming recognize: %w", err)
	}

	return stream, nil
}

func (sc *speechClient) Close() error {
	return sc.c.Close()
}
