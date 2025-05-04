package speech_wrapper

import "cloud.google.com/go/speech/apiv1/speechpb"

type streamingClient struct {
	stream speechpb.Speech_StreamingRecognizeClient
}

func (s *streamingClient) Send() (*speechpb.StreamingRecognizeRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s *streamingClient) Recv() (*speechpb.StreamingRecognizeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *streamingClient) CloseSend() error {
	//TODO implement me
	panic("implement me")
}
