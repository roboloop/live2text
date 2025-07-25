package btt_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"live2text/internal/services/btt"
	"live2text/internal/services/btt/client"
	"live2text/internal/services/btt/client/trigger"
	"live2text/internal/services/btt/storage"
	"live2text/internal/services/btt/tmpl"
	"live2text/internal/services/recognition"
	"live2text/internal/utils/logger"
)

func TestToggleListening(t *testing.T) {
	t.Parallel()

	t.Run("cannot get task id", func(t *testing.T) {
		t.Parallel()

		lc := setupListeningComponent(
			t,
			func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("", errors.New("something happened"))
			},
			nil,
		)

		err := lc.ToggleListening(t.Context())
		require.Error(t, err)
		require.ErrorContains(t, err, "cannot get task id")
	})

	t.Run("start listening", func(t *testing.T) {
		t.Parallel()

		lc := setupListeningComponent(
			t,
			func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("", nil)

				// dirty hack to prevent next code execution
				dc.SelectedDeviceMock.Return("", errors.New("must exit"))
			},
			nil,
		)

		err := lc.ToggleListening(t.Context())
		require.Error(t, err)
		require.ErrorContains(t, err, "must exit")
	})

	t.Run("stop listening", func(t *testing.T) {
		t.Parallel()

		lc := setupListeningComponent(
			t,
			func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				getValueCalls := 0
				s.GetValueMock.Set(func(ctx context.Context, key storage.Key) (string, error) {
					getValueCalls += 1
					switch getValueCalls {
					case 1:
						return "foo", nil
					default:
						// dirty hack to prevent next code execution
						return "", errors.New("must exit")
					}
				})
			},
			nil,
		)

		err := lc.ToggleListening(t.Context())
		require.Error(t, err)
		require.ErrorContains(t, err, "must exit")
	})
}

func TestStartListening(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock)
		expectErr  string
	}{
		{
			name: "cannot get selected device",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Return("", errors.New("something happened"))
			},
			expectErr: "cannot get selected device",
		},
		{
			name: "cannot get selected language",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Return("foo", nil)
				lc.SelectedLanguageMock.Return("", errors.New("something happened"))
			},
			expectErr: "cannot get selected language",
		},
		{
			name: "cannot start recognition",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Return("foo", nil)
				lc.SelectedLanguageMock.Return("bar", nil)
				rg.StartMock.Return("", "", errors.New("something happened"))
			},
			expectErr: "cannot start recognition",
		},
		{
			name: "cannot set task id",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Return("foo", nil)
				lc.SelectedLanguageMock.Return("bar", nil)
				rg.StartMock.Return("baz", "qux", nil)
				s.SetValueMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot set task id",
		},
		{
			name: "cannot update app",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Return("foo", nil)
				lc.SelectedLanguageMock.Return("bar", nil)
				rg.StartMock.Return("baz", "qux", nil)
				s.SetValueMock.Return(nil)
				r.ListenSocketMock.Return("quux")
				c.UpdateTriggerMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot update app",
		},
		{
			name: "cannot update clean view app",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Return("foo", nil)
				lc.SelectedLanguageMock.Return("bar", nil)
				rg.StartMock.Return("baz", "qux", nil)
				s.SetValueMock.Return(nil)
				r.ListenSocketMock.Return("quux")
				updateTriggerCalls := 0
				c.UpdateTriggerMock.Set(
					func(ctx context.Context, title trigger.Title, patch trigger.Trigger) error {
						updateTriggerCalls += 1
						switch updateTriggerCalls {
						case 1:
							return nil
						default:
							return errors.New("something happened")
						}
					},
				)
			},
			expectErr: "cannot update clean view app",
		},
		{
			name: "cannot enable clean mode",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Return("foo", nil)
				lc.SelectedLanguageMock.Return("bar", nil)
				rg.StartMock.Return("baz", "qux", nil)
				s.SetValueMock.Return(nil)
				r.ListenSocketMock.Return("quux")
				c.UpdateTriggerMock.Times(2).Return(nil)
				vmc.EnableCleanModeMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot enable clean mode",
		},
		{
			name: "cannot show floating",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Return("foo", nil)
				lc.SelectedLanguageMock.Return("bar", nil)
				rg.StartMock.Return("baz", "qux", nil)
				s.SetValueMock.Return(nil)
				r.ListenSocketMock.Return("quux")
				c.UpdateTriggerMock.Times(2).Return(nil)
				vmc.EnableCleanModeMock.Return(nil)
				fc.ShowFloatingMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot show floating",
		},
		{
			name: "cannot show clipboard",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Return("foo", nil)
				lc.SelectedLanguageMock.Return("bar", nil)
				rg.StartMock.Return("baz", "qux", nil)
				s.SetValueMock.Return(nil)
				r.ListenSocketMock.Return("quux")
				c.UpdateTriggerMock.Times(2).Return(nil)
				vmc.EnableCleanModeMock.Return(nil)
				fc.ShowFloatingMock.Return(nil)
				cc.ShowClipboardMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot show clipboard",
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				dc.SelectedDeviceMock.Expect(minimock.AnyContext).Return("foo", nil)
				lc.SelectedLanguageMock.Expect(minimock.AnyContext).Return("bar", nil)
				rg.StartMock.Expect(minimock.AnyContext, "foo", "bar").Return("baz", "qux", nil)
				s.SetValueMock.Expect(minimock.AnyContext, "LIVE2TEXT_TASK_ID", "baz").Return(nil)
				r.ListenSocketMock.Expect("qux").Return("quux")
				updateTriggerCalls := 0
				c.UpdateTriggerMock.Inspect(func(ctx context.Context, title trigger.Title, patch trigger.Trigger) {
					updateTriggerCalls += 1
					switch updateTriggerCalls {
					case 1:
						require.EqualValues(t, "App", title)
					case 2:
						require.EqualValues(t, "Clean View App", title)
					}
				}).Times(2).Return(nil)
				vmc.EnableCleanModeMock.Expect(minimock.AnyContext).Return(nil)
				fc.ShowFloatingMock.Expect(minimock.AnyContext).Return(nil)
				cc.ShowClipboardMock.Expect(minimock.AnyContext).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lc := setupListeningComponent(t, tt.setupMocks, nil)

			err := lc.StartListening(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestStopListening(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock)
		expectErr  string
	}{
		{
			name: "cannot get task id",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("", errors.New("something happened"))
			},
			expectErr: "cannot get task id",
		},
		{
			name: "cannot empty task id",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("foo", nil)
				rg.StopMock.Return(nil)
				s.SetValueMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot empty task id",
		},
		{
			name: "cannot update app",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("foo", nil)
				rg.StopMock.Return(nil)
				s.SetValueMock.Return(nil)
				r.AppPlaceholderMock.Return("")
				c.UpdateTriggerMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot update app",
		},
		{
			name: "cannot update clean view app",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("foo", nil)
				rg.StopMock.Return(nil)
				s.SetValueMock.Return(nil)
				r.AppPlaceholderMock.Return("")
				updateTriggerCalls := 0
				c.UpdateTriggerMock.Set(
					func(ctx context.Context, title trigger.Title, patch trigger.Trigger) error {
						updateTriggerCalls += 1
						switch updateTriggerCalls {
						case 1:
							return nil
						default:
							return errors.New("something happened")
						}
					},
				)
			},
			expectErr: "cannot update clean view app",
		},
		{
			name: "cannot disable clean mode",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("foo", nil)
				rg.StopMock.Return(nil)
				s.SetValueMock.Return(nil)
				r.AppPlaceholderMock.Return("")
				c.UpdateTriggerMock.Times(2).Return(nil)
				vmc.DisableCleanViewMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot disable clean mode",
		},
		{
			name: "cannot hide floating",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("foo", nil)
				rg.StopMock.Return(nil)
				s.SetValueMock.Return(nil)
				r.AppPlaceholderMock.Return("")
				c.UpdateTriggerMock.Times(2).Return(nil)
				vmc.DisableCleanViewMock.Return(nil)
				fc.HideFloatingMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot hide floating",
		},
		{
			name: "cannot hide clipboard",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("foo", nil)
				rg.StopMock.Return(nil)
				s.SetValueMock.Return(nil)
				r.AppPlaceholderMock.Return("")
				c.UpdateTriggerMock.Times(2).Return(nil)
				vmc.DisableCleanViewMock.Return(nil)
				fc.HideFloatingMock.Return(nil)
				cc.HideClipboardMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot hide clipboard",
		},
		{
			name: "happy path",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Expect(minimock.AnyContext, "LIVE2TEXT_TASK_ID").Return("foo", nil)
				rg.StopMock.Expect(minimock.AnyContext, "foo").Return(nil)
				s.SetValueMock.Return(nil)
				r.AppPlaceholderMock.Expect().Return("")
				updateTriggerCalls := 0
				c.UpdateTriggerMock.Inspect(func(ctx context.Context, title trigger.Title, patch trigger.Trigger) {
					updateTriggerCalls += 1
					switch updateTriggerCalls {
					case 1:
						require.EqualValues(t, "App", title)
					case 2:
						require.EqualValues(t, "Clean View App", title)
					}
				}).Times(2).Return(nil)
				vmc.DisableCleanViewMock.Expect(minimock.AnyContext).Return(nil)
				fc.HideFloatingMock.Expect(minimock.AnyContext).Return(nil)
				cc.HideClipboardMock.Expect(minimock.AnyContext).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lc := setupListeningComponent(t, tt.setupMocks, nil)

			err := lc.StopListening(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}

	t.Run("cannot stop task logged", func(t *testing.T) {
		t.Parallel()

		l, h := logger.NewCaptureLogger()
		lc := setupListeningComponent(
			t,
			func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("foo", nil)
				rg.StopMock.Return(errors.New("something happened"))

				// dirty hack to prevent next code execution
				s.SetValueMock.Return(errors.New("must exit"))
			},
			l,
		)

		_ = lc.StopListening(t.Context())

		log, ok := h.GetLog("Cannot stop task")
		require.Truef(t, ok, "cannot find the log entry")
		errAttr, ok := log.GetAttr("error")
		require.Truef(t, ok, "cannot find the error attribute")
		require.Implements(t, (*error)(nil), errAttr)
		require.ErrorContains(t, errAttr.(error), "something happened")
	})
}

func TestIsRunning(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock)
		expectRunning bool
		expectErr     string
	}{
		{
			name: "cannot get task id",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Return("", errors.New("something happened"))
			},
			expectErr: "cannot get task id",
		},
		{
			name: "task is running",
			setupMocks: func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock) {
				s.GetValueMock.Expect(minimock.AnyContext, "LIVE2TEXT_TASK_ID").Return("foo", nil)
				rg.HasMock.Expect("foo").Return(true)
			},
			expectRunning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lc := setupListeningComponent(t, tt.setupMocks, nil)

			isRunning, err := lc.IsRunning(t.Context())
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectRunning, isRunning)
		})
	}
}

func setupListeningComponent(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, rg *recognition.RecognitionMock, c *client.ClientMock, s *storage.StorageMock, r *tmpl.RendererMock, dc *btt.DeviceComponentMock, lc *btt.LanguageComponentMock, vmc *btt.ViewModeComponentMock, fc *btt.FloatingComponentMock, cc *btt.ClipboardComponentMock),
	l *slog.Logger,
) btt.ListeningComponent {
	mc := minimock.NewController(t)
	rg := recognition.NewRecognitionMock(mc)
	c := client.NewClientMock(mc)
	s := storage.NewStorageMock(mc)
	r := tmpl.NewRendererMock(mc)
	dc := btt.NewDeviceComponentMock(mc)
	lc := btt.NewLanguageComponentMock(mc)
	vmc := btt.NewViewModeComponentMock(mc)
	fc := btt.NewFloatingComponentMock(mc)
	cc := btt.NewClipboardComponentMock(mc)

	if setupMocks != nil {
		setupMocks(mc, rg, c, s, r, dc, lc, vmc, fc, cc)
	}
	if l == nil {
		l = logger.NilLogger
	}

	return btt.NewListeningComponent(l, rg, c, s, r, dc, lc, vmc, fc, cc)
}
