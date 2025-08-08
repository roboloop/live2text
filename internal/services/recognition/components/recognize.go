package components

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/roboloop/live2text/internal/services/metrics"
	speechwrapper "github.com/roboloop/live2text/internal/services/speech_wrapper"
	"github.com/roboloop/live2text/internal/utils/logger"
)

type RecognizeParameters struct {
	Channels   int
	SampleRate int
	Language   string
}

type Recognized struct {
	Transcript string
	IsFinal    bool
	EndTime    time.Duration
}

//go:generate minimock -g -i RecognizerComponent -s _mock.go -o .

type RecognizerComponent interface {
	Recognize(ctx context.Context, parameters RecognizeParameters, input <-chan []int16, output chan<- Recognized) error
}

type recognizerComponent struct {
	logger         *slog.Logger
	metrics        metrics.Metrics
	speechClient   speechwrapper.SpeechClient
	streamDuration time.Duration
}

func NewRecognizerComponent(
	logger *slog.Logger,
	metrics metrics.Metrics,
	speechClient speechwrapper.SpeechClient,
) RecognizerComponent {
	return &recognizerComponent{
		logger:         logger,
		metrics:        metrics,
		speechClient:   speechClient,
		streamDuration: 5 * time.Minute,
	}
}

// Recognize runs the recognition. The `output` must be closed after the function is completed.
func (r *recognizerComponent) Recognize(
	ctx context.Context,
	parameters RecognizeParameters,
	input <-chan []int16,
	output chan<- Recognized,
) error {
	for {
		if errRecognize := r.recognize(ctx, parameters, input, output); errRecognize != nil {
			return errRecognize
		}
	}
}

func (r *recognizerComponent) recognize(
	ctx context.Context,
	parameters RecognizeParameters,
	input <-chan []int16,
	output chan<- Recognized,
) error {
	r.logger.InfoContext(ctx, "New streaming recognize request")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := r.speechClient.StreamingRecognize(ctx)
	if err != nil {
		return fmt.Errorf("cannot create a streaming connection: %w", err)
	}
	defer func() {
		if errClose := stream.CloseSend(); errClose != nil {
			r.logger.ErrorContext(ctx, "Error closing stream", "error", errClose)
		}
	}()

	if err = stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: int32(parameters.SampleRate), //nolint:gosec // the value always small enough
					LanguageCode:    parameters.Language,
					MaxAlternatives: 1,
				},
				InterimResults: true,
			},
		},
	}); err != nil {
		return fmt.Errorf("cannot send a config request: %w", err)
	}

	r.metrics.AddConnectionsToGoogleSpeech(1)
	wg := sync.WaitGroup{}

	sendErrCh := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(sendErrCh)

		if sendErr := r.send(ctx, stream, parameters, input); sendErr != nil {
			r.logger.Log(ctx, logger.ResolveLevel(sendErr), "Shutting down send loop", "error", sendErr)

			sendErrCh <- sendErr
		}
	}()

	receiveErrCh := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(receiveErrCh)

		if receiveErr := r.receive(ctx, stream, output); receiveErrCh != nil {
			r.logger.Log(ctx, logger.ResolveLevel(receiveErr), "Shutting down receive loop", "error", receiveErr)

			receiveErrCh <- receiveErr
		}
	}()

	var terminationErr error
	select {
	case <-ctx.Done():
		r.logger.InfoContext(ctx, "Stream exiting by parent context", "error", ctx.Err())
		terminationErr = ctx.Err()
	case terminationErr = <-sendErrCh:
		r.logger.ErrorContext(ctx, "Stream exiting by sent error...", "error", terminationErr)
	case terminationErr = <-receiveErrCh:
		r.logger.ErrorContext(ctx, "Stream exiting by received error...", "error", terminationErr)
	case <-time.After(r.streamDuration):
		r.logger.InfoContext(ctx, "Stream exiting by timer (restart)...")
	}

	cancel()
	wg.Wait()

	// drain them
	<-sendErrCh
	<-receiveErrCh

	return terminationErr
}

func (r *recognizerComponent) send(
	ctx context.Context,
	stream speechwrapper.StreamingRecognizeClient,
	parameters RecognizeParameters,
	input <-chan []int16,
) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case buffer, ok := <-input:
			if !ok {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				return errors.New("input channel closed unexpectedly")
			}
			content := encodeInt16ToBytes(buffer)
			durationMs := snippetDuration(buffer, parameters.Channels, parameters.SampleRate)

			if err = stream.Send(&speechpb.StreamingRecognizeRequest{
				StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
					AudioContent: content,
				},
			}); err != nil {
				return fmt.Errorf("cannot send audio: %w", err)
			}

			r.metrics.AddBytesSentToGoogleSpeech(len(content))
			r.metrics.AddMillisecondsSentToGoogleSpeech(durationMs)
		}
	}
}

//nolint:gocognit // TODO: heavy error handling logic
func (r *recognizerComponent) receive(
	ctx context.Context,
	stream speechwrapper.StreamingRecognizeClient,
	recognizedCh chan<- Recognized,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			if st, ok := status.FromError(err); ok && st.Code() == codes.Canceled {
				return ctx.Err()
			}

			return fmt.Errorf("cannot receive results: %w", err)
		}
		if grpcError := resp.GetError(); grpcError != nil {
			st := status.FromProto(grpcError)
			// Workaround while the API doesn't give a more informative error.
			if st.Code() == codes.InvalidArgument || st.Code() == codes.OutOfRange {
				r.logger.WarnContext(ctx, "Speech recognition request exceeded a limit of 300 seconds")
			}
			if st.Code() == codes.Canceled {
				return ctx.Err()
			}
			return fmt.Errorf("cannot recognize audio: %w", st.Err())
		}

		if resp.GetSpeechEventType() != speechpb.StreamingRecognizeResponse_SPEECH_EVENT_UNSPECIFIED {
			r.logger.WarnContext(ctx, "Unusual event", "resp", resp)
			continue
		}

		if len(resp.GetResults()) == 0 {
			continue
		}

		result := resp.GetResults()[0]
		transcript := strings.TrimSpace(result.GetAlternatives()[0].GetTranscript())
		recognizedCh <- Recognized{
			Transcript: transcript,
			IsFinal:    result.GetIsFinal(),
			EndTime:    result.GetResultEndTime().AsDuration(),
		}
	}
}

func encodeInt16ToBytes(data []int16) []byte {
	content := make([]byte, 0, len(data)*2)
	for _, value := range data {
		lowByte := byte(value & math.MaxUint8)
		highByte := byte((value >> 8) & math.MaxUint8)
		content = append(content, lowByte, highByte)
	}
	return content
}

func snippetDuration(data []int16, channels, sampleRate int) int {
	totalSamples := len(data) / (16 / 8 * channels)
	durationMs := totalSamples * 1000 / sampleRate

	return durationMs
}
