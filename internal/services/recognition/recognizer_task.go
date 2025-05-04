package recognition

import (
	"cloud.google.com/go/speech/apiv1/speechpb"
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"live2text/internal/background"
	"live2text/internal/services/audio"
	"live2text/internal/services/burner"
	"live2text/internal/services/metrics"
	"live2text/internal/services/recognition/console"
	"live2text/internal/services/recognition/subs"
	"live2text/internal/services/speech_wrapper"
	"live2text/internal/utils"
	"log"
	"log/slog"
	"math"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type recognized struct {
	transcript string
	isFinal    bool
	endTime    time.Duration
}

type RecognizeTask struct {
	logger        *slog.Logger
	metrics       metrics.Metrics
	audio         audio.Audio
	burner        burner.Burner
	socketManager *background.SocketManager
	speechClient  speech_wrapper.Client

	id         string
	device     string
	language   string
	socketPath string

	streamDuration time.Duration
	subs           *subs.Writer
}

func NewRecognizeTask(logger *slog.Logger, metrics metrics.Metrics, audio audio.Audio, burner burner.Burner, socketManager *background.SocketManager, speechClient speech_wrapper.Client, device, language, socketPath string) *RecognizeTask {
	return &RecognizeTask{
		logger:        logger,
		metrics:       metrics,
		audio:         audio,
		burner:        burner,
		socketManager: socketManager,
		speechClient:  speechClient,

		id:         fmt.Sprintf("%d", rand.Int()%100000),
		device:     device,
		language:   language,
		socketPath: socketPath,

		streamDuration: 5 * time.Minute,
		subs:           subs.NewWriter(2, 80),
	}
}

func (rt *RecognizeTask) Run(ctx context.Context) error {
	deviceListener, err := rt.audio.ListenDevice(ctx, rt.device)
	if err != nil {
		return fmt.Errorf("cannot listen device: %w", err)
	}

	audioBroadcaster := utils.Broadcaster(ctx, rt.logger, deviceListener.Ch, 2)

	g, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	// Listen device
	wg.Add(1)
	g.Go(func() error {
		defer wg.Done()
		return deviceListener.Listen(ctx)
	})

	// Burn listening part on device
	wg.Add(1)
	g.Go(func() error {
		defer wg.Done()
		return rt.burnContent(ctx, audioBroadcaster[0], deviceListener.Channels, deviceListener.SampleRate)
	})

	recognizedCh := make(chan recognized, 1024)
	// Recognize listening part
	wg.Add(1)
	g.Go(func() error {
		defer wg.Done()
		defer close(recognizedCh)
		for {
			if streamErr := rt.stream(ctx, audioBroadcaster[1], deviceListener.SampleRate, recognizedCh); streamErr != nil {
				return streamErr
			}
		}
	})

	subsBroadcaster := utils.Broadcaster(ctx, rt.logger, recognizedCh, 2)

	// Store recognized part in internal memory
	wg.Add(1)
	g.Go(func() error {
		defer wg.Done()

		return rt.storeSubs(ctx, subsBroadcaster[0])
	})

	// Print recognized part to stdout
	wg.Add(1)
	g.Go(func() error {
		defer wg.Done()

		return rt.printSubs(ctx, subsBroadcaster[1])
	})

	// Print to socket manager
	wg.Add(1)
	g.Go(func() error {
		defer wg.Done()

		return rt.socketManager.Listen(rt.socketPath, func(conn net.Conn) {
			defer conn.Close()
			conn.Write([]byte(rt.subs.Format()))
		})
	})

	if err = g.Wait(); err != nil {
		slog.ErrorContext(ctx, "Recognition failed", "error", err)
	}

	wg.Wait()

	return err
}

func (rt *RecognizeTask) burnContent(ctx context.Context, ch <-chan []int16, channels int, sampleRate int) error {
	name := time.Now().Format("01.02.06 15_04_05 output.wav")
	file, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer file.Close()

	if err = rt.burner.Burn(ctx, file, ch, channels, sampleRate); err != nil {
		return fmt.Errorf("cannot burn: %w", err)
	}

	return err
}

func (rt *RecognizeTask) stream(ctx context.Context, ch <-chan []int16, sampleRate int, recognizedCh chan<- recognized) error {
	rt.logger.InfoContext(ctx, "New streaming recognize request")

	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := rt.speechClient.StreamingRecognize(streamCtx)
	if err != nil {
		return fmt.Errorf("could not streaming recognize: %w", err)
	}
	defer stream.CloseSend()
	if err = stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: int32(sampleRate),
					LanguageCode:    rt.language,
					MaxAlternatives: 1,
				},
				InterimResults: true,
			},
		},
	}); err != nil {
		return fmt.Errorf("could not send config request: %w", err)
	}

	g, errCtx := errgroup.WithContext(streamCtx)
	g.Go(func() error {
		return rt.streamContent(errCtx, stream, ch)
	})
	g.Go(func() error {
		return rt.readRecognized(errCtx, stream, recognizedCh)
	})

	errCh := make(chan error, 1)
	go func() {
		if groupErr := g.Wait(); groupErr != nil {
			errCh <- groupErr
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		rt.logger.InfoContext(ctx, "Stream exiting by parent ctx...", "error", ctx.Err())
		return ctx.Err()
	case err = <-errCh:
		rt.logger.InfoContext(ctx, "Stream exiting by error...", "error", err)
		return err
	case <-time.NewTimer(rt.streamDuration).C:
		rt.logger.InfoContext(ctx, "Stream exiting by timer (restart)...")
		return nil
	}
}

func (rt *RecognizeTask) streamContent(ctx context.Context, stream speechpb.Speech_StreamingRecognizeClient, ch <-chan []int16) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			rt.logger.InfoContext(ctx, "[Stream Content] Shutting down...")
			return ctx.Err()
		case buffer := <-ch:
			//rt.logger.InfoContext(ctx, "Receiving the buf", "size", len(buffer))
			content := make([]byte, 0, len(buffer)*2)
			for _, value := range buffer {
				lowByte := byte(value & math.MaxUint8)
				highByte := byte((value >> 8) & math.MaxUint8)
				content = append(content, lowByte, highByte)
			}

			if err = stream.Send(&speechpb.StreamingRecognizeRequest{
				StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
					AudioContent: content,
				},
			}); err != nil {
				return fmt.Errorf("could not send audio: %w", err)
			}

			rt.metrics.AddBytesSentToGoogleSpeech(len(content))
		}
	}
}

func (rt *RecognizeTask) readRecognized(ctx context.Context, stream speechpb.Speech_StreamingRecognizeClient, recognizedCh chan<- recognized) error {
	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "[Read recognized] Shutting down...")
			return ctx.Err()
		default:
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("cannot stream results: %w", err)
		}
		if err := resp.Error; err != nil {
			// Workaround while the API doesn't give a more informative error.
			if err.Code == 3 || err.Code == 11 {
				log.Print("WARNING: Speech recognition request exceeded limit of 60 seconds.")
			}
			return fmt.Errorf("could not recognize: %v", err)
		}

		if resp.GetSpeechEventType() != speechpb.StreamingRecognizeResponse_SPEECH_EVENT_UNSPECIFIED {
			rt.logger.WarnContext(ctx, "Unusual event", "resp", resp)
			continue
		}

		if len(resp.Results) == 0 {
			continue
		}

		result := resp.Results[0]
		transcript := strings.TrimSpace(result.Alternatives[0].Transcript)
		recognizedCh <- recognized{
			transcript: transcript,
			isFinal:    result.IsFinal,
			endTime:    result.ResultEndTime.AsDuration(),
		}
	}

	return nil
}

func (rt *RecognizeTask) storeSubs(_ context.Context, sectionsCh <-chan recognized) error {
	for s := range sectionsCh {
		rt.subs.AddSection(s.transcript, s.isFinal)
	}
	return nil
}

func (rt *RecognizeTask) printSubs(_ context.Context, recognizedCh <-chan recognized) error {
	fd := os.NewFile(uintptr(3), "fd3")
	if fd == nil {
		return errors.New("could not open fd")
	}
	defer fd.Close()

	var lastWasFinal bool
	writer := console.NewWriter(fd)
	for r := range recognizedCh {
		text := fmt.Sprintf("%s: %s", r.endTime.Truncate(time.Second).String(), r.transcript)
		if r.isFinal {
			writer.PrintSuccess(text)
			lastWasFinal = true
		} else {
			writer.PrintFail(text)
			lastWasFinal = false
		}
	}
	if !lastWasFinal {
		writer.PrintNewLine()
	}

	return nil
}

func (rt *RecognizeTask) Subs() string {
	return rt.subs.Format()
}

func (rt *RecognizeTask) Close() error {
	return nil
}
