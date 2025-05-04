package speech_wrapper

import (
	"cloud.google.com/go/speech/apiv1/speechpb"
	"context"
)

type Client interface {
	StreamingRecognize(ctx context.Context) (speechpb.Speech_StreamingRecognizeClient, error)
	Close() error
}

//type StreamingClient interface {
//	Send() (*speechpb.StreamingRecognizeRequest, error)
//	Recv() (*speechpb.StreamingRecognizeResponse, error)
//	CloseSend() error
//}
