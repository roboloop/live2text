package components_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"cloud.google.com/go/speech/apiv1/speechpb"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/roboloop/live2text/internal/services/metrics"
	"github.com/roboloop/live2text/internal/services/recognition/components"
	speechwrapper "github.com/roboloop/live2text/internal/services/speech_wrapper"
	"github.com/roboloop/live2text/internal/utils/logger"
)

// TODO: more tests

func TestRecognize(t *testing.T) {
	t.Parallel()

	defaultParameters := components.RecognizeParameters{
		Channels:   1,
		SampleRate: 24000,
		Language:   "en-US",
	}

	t.Run("cannot create a streaming connection", func(t *testing.T) {
		t.Parallel()

		recognition := setupRecognition(
			t,
			func(mc *minimock.Controller, m *metrics.MetricsMock, sc *speechwrapper.SpeechClientMock) {
				sc.StreamingRecognizeMock.Return(nil, errors.New("something happened"))
			},
			nil,
		)

		err := recognition.Recognize(t.Context(), defaultParameters, nil, nil)

		require.Error(t, err)
		require.ErrorContains(t, err, "cannot create a streaming connection")
		require.ErrorContains(t, err, "something happened")
	})

	t.Run("cannot send a config request", func(t *testing.T) {
		t.Parallel()

		recognition := setupRecognition(
			t,
			func(mc *minimock.Controller, m *metrics.MetricsMock, sc *speechwrapper.SpeechClientMock) {
				src := speechwrapper.NewStreamingRecognizeClientMock(mc)
				src.CloseSendMock.Return(nil)
				src.SendMock.Return(errors.New("something happened"))

				sc.StreamingRecognizeMock.Return(src, nil)
			},
			nil,
		)

		err := recognition.Recognize(t.Context(), defaultParameters, nil, nil)

		require.Error(t, err)
		require.ErrorContains(t, err, "cannot send a config request")
		require.ErrorContains(t, err, "something happened")
	})

	t.Run("cannot send audio", func(t *testing.T) {
		t.Parallel()

		l, h := logger.NewCaptureLogger()
		recognition := setupRecognition(
			t,
			func(mc *minimock.Controller, m *metrics.MetricsMock, sc *speechwrapper.SpeechClientMock) {
				src := speechwrapper.NewStreamingRecognizeClientMock(mc)
				src.CloseSendMock.Return(nil)
				sendCalls := 0
				src.SendMock.Set(func(*speechpb.StreamingRecognizeRequest) error {
					sendCalls += 1
					switch sendCalls {
					case 1:
						return nil
					default:
						return errors.New("something happened")
					}
				})
				m.AddConnectionsToGoogleSpeechMock.Expect(1).Return()
				src.RecvMock.Set(func() (*speechpb.StreamingRecognizeResponse, error) {
					time.Sleep(10 * time.Millisecond)
					return nil, errors.New("abort further execution")
				})

				sc.StreamingRecognizeMock.Return(src, nil)
			},
			l,
		)

		input := make(chan []int16, 1)
		defer close(input)
		input <- []int16{0, 1, 2, 3}
		output := make(chan components.Recognized, 1)
		defer close(output)

		err := recognition.Recognize(t.Context(), defaultParameters, input, output)

		require.Error(t, err)
		require.ErrorContains(t, err, "cannot send audio")
		require.ErrorContains(t, err, "something happened")

		requireLogError(t, h, slog.LevelError, "Stream exiting by sent error...", "something happened")
	})

	t.Run("cannot receive a response", func(t *testing.T) {
		t.Parallel()

		l, h := logger.NewCaptureLogger()
		recognition := setupRecognition(
			t,
			func(mc *minimock.Controller, m *metrics.MetricsMock, sc *speechwrapper.SpeechClientMock) {
				src := speechwrapper.NewStreamingRecognizeClientMock(mc)
				sc.StreamingRecognizeMock.Return(src, nil)
				src.CloseSendMock.Return(nil)
				sendCalls := 0
				src.SendMock.Set(func(req *speechpb.StreamingRecognizeRequest) error {
					sendCalls += 1
					switch sendCalls {
					case 1:
						return nil
					default:
						time.Sleep(10 * time.Millisecond)
						return errors.New("abort further execution")
					}
				})
				m.AddConnectionsToGoogleSpeechMock.Expect(1).Return()
				src.RecvMock.Set(func() (*speechpb.StreamingRecognizeResponse, error) {
					return nil, errors.New("something happened")
				})
			},
			l,
		)

		input := make(chan []int16, 1)
		input <- []int16{0, 1, 2, 3}
		err := recognition.Recognize(t.Context(), defaultParameters, input, nil)

		require.Error(t, err)
		require.ErrorContains(t, err, "cannot receive results")
		require.ErrorContains(t, err, "something happened")

		requireLogError(t, h, slog.LevelError, "Stream exiting by received error...", "something happened")
	})

	t.Run("response is error", func(t *testing.T) {
		t.Parallel()

		l, h := logger.NewCaptureLogger()
		recognition := setupRecognition(
			t,
			func(mc *minimock.Controller, m *metrics.MetricsMock, sc *speechwrapper.SpeechClientMock) {
				src := speechwrapper.NewStreamingRecognizeClientMock(mc)
				sc.StreamingRecognizeMock.Return(src, nil)
				src.CloseSendMock.Return(nil)
				sendCalls := 0
				src.SendMock.Set(func(req *speechpb.StreamingRecognizeRequest) error {
					sendCalls += 1
					switch sendCalls {
					case 1:
						return nil
					default:
						time.Sleep(10 * time.Millisecond)
						return errors.New("abort further execution")
					}
				})
				m.AddConnectionsToGoogleSpeechMock.Expect(1).Return()
				src.RecvMock.Set(func() (*speechpb.StreamingRecognizeResponse, error) {
					return &speechpb.StreamingRecognizeResponse{
						Error: status.New(codes.InvalidArgument, "something happened").Proto(),
					}, nil
				})
			},
			l,
		)

		input := make(chan []int16, 1)
		input <- []int16{0, 1, 2, 3}
		err := recognition.Recognize(t.Context(), defaultParameters, input, nil)

		require.Error(t, err)
		require.ErrorContains(t, err, "cannot recognize audio")
		require.ErrorContains(t, err, "something happened")

		requireLogError(t, h, slog.LevelError, "Stream exiting by received error...", "something happened")
	})

	t.Run("recognize one audio snippet", func(t *testing.T) {
		t.Parallel()

		l, _ := logger.NewCaptureLogger()
		recognition := setupRecognition(
			t,
			func(mc *minimock.Controller, m *metrics.MetricsMock, sc *speechwrapper.SpeechClientMock) {
				src := speechwrapper.NewStreamingRecognizeClientMock(mc)
				sc.StreamingRecognizeMock.Return(src, nil)
				src.CloseSendMock.Return(nil)
				src.SendMock.Return(nil)
				m.AddConnectionsToGoogleSpeechMock.Return()
				m.AddBytesSentToGoogleSpeechMock.Return()
				m.AddMillisecondsSentToGoogleSpeechMock.Return()

				recvCalls := 0
				src.RecvMock.Set(func() (*speechpb.StreamingRecognizeResponse, error) {
					recvCalls += 1
					switch recvCalls {
					case 1:
						return &speechpb.StreamingRecognizeResponse{
							Error:           nil,
							SpeechEventType: speechpb.StreamingRecognizeResponse_SPEECH_EVENT_UNSPECIFIED,
							Results: []*speechpb.StreamingRecognitionResult{
								{
									Alternatives: []*speechpb.SpeechRecognitionAlternative{
										{
											Transcript: " sample text ",
										},
									},
									IsFinal:       true,
									ResultEndTime: durationpb.New(100 * time.Millisecond),
								},
							},
						}, nil
					default:
						return nil, io.EOF
					}
				})
			},
			l,
		)

		input := make(chan []int16, 1)
		input <- []int16{255, 256, 257, 258}
		defer close(input)
		output := make(chan components.Recognized, 1)
		defer close(output)

		ctx, cancel := context.WithCancel(t.Context())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		err := recognition.Recognize(ctx, defaultParameters, input, output)
		recognized := <-output

		require.Error(t, err)
		require.ErrorIs(t, err, context.Canceled)
		require.Equal(t, components.Recognized{
			Transcript: "sample text",
			IsFinal:    true,
			EndTime:    100 * time.Millisecond,
		}, recognized)
	})
}

func setupRecognition(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, m *metrics.MetricsMock, sc *speechwrapper.SpeechClientMock),
	l *slog.Logger,
) components.RecognizerComponent {
	mc := minimock.NewController(t)
	m := metrics.NewMetricsMock(mc)
	sc := speechwrapper.NewSpeechClientMock(mc)

	if setupMocks != nil {
		setupMocks(mc, m, sc)
	}
	if l == nil {
		l = logger.NilLogger
	}

	return components.NewRecognizerComponent(l, m, sc)
}

func requireLogError(t *testing.T, h *logger.CaptureHandler, level slog.Level, msg, errorMsg string) {
	t.Helper()

	logEntry, ok := h.GetLog(msg)
	require.Truef(t, ok, "cannot find the log entry")
	require.Equal(t, level, logEntry.Level)

	errAttr, ok := logEntry.GetAttr("error")
	require.Truef(t, ok, "cannot find the error attribute")
	require.Implements(t, (*error)(nil), errAttr)
	require.ErrorContains(t, errAttr.(error), errorMsg)
}
