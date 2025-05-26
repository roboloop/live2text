package speechwrapper

import (
	"context"

	"cloud.google.com/go/speech/apiv1/speechpb"
)

type Client interface {
	StreamingRecognize(ctx context.Context) (speechpb.Speech_StreamingRecognizeClient, error)
	Close() error
}
