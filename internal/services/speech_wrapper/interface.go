package speechwrapper

import (
	"context"

	"cloud.google.com/go/speech/apiv1/speechpb"
)

//go:generate minimock -g -i SpeechClient,StreamingRecognizeClient -s _mock.go -o .

type SpeechClient interface {
	StreamingRecognize(ctx context.Context) (StreamingRecognizeClient, error)
	Close() error
}

type StreamingRecognizeClient interface {
	Send(*speechpb.StreamingRecognizeRequest) error
	Recv() (*speechpb.StreamingRecognizeResponse, error)
	CloseSend() error
}
