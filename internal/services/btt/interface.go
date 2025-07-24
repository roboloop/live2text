package btt

import (
	"context"

	"live2text/internal/services/btt/client/trigger"
	"live2text/internal/services/btt/storage"
)

//go:generate minimock -g -i Btt,InitializingComponent,ListeningComponent,DeviceComponent,LanguageComponent,ViewModeComponent,FloatingComponent,SettingsComponent -s _mock.go -o .

type Btt interface {
	InitializingComponent
	ListeningComponent
	DeviceComponent
	LanguageComponent
	ViewModeComponent
	FloatingComponent
}

type InitializingComponent interface {
	Initialize(ctx context.Context) error
	Clear(ctx context.Context) error
}

type ListeningComponent interface {
	ToggleListening(ctx context.Context) error
	StartListening(ctx context.Context) error
	StopListening(ctx context.Context) error
	IsRunning(ctx context.Context) (bool, error)
}

type DeviceComponent interface {
	LoadDevices(ctx context.Context) error
	SelectDevice(ctx context.Context, device string) error
	SelectedDevice(ctx context.Context) (string, error)
}

type LanguageComponent interface {
	SelectLanguage(ctx context.Context, language string) error
	SelectedLanguage(ctx context.Context) (string, error)
}

type ViewModeComponent interface {
	SelectViewMode(ctx context.Context, viewMode ViewMode) error
	SelectedViewMode(ctx context.Context) (ViewMode, error)
	EnableCleanMode(ctx context.Context) error
	DisableCleanView(ctx context.Context) error
}

type FloatingComponent interface {
	SelectFloating(ctx context.Context, floating Floating) error
	SelectedFloating(ctx context.Context) (Floating, error)
	ShowFloating(ctx context.Context) error
	HideFloating(ctx context.Context) error
	FloatingPage() string
	StreamText(ctx context.Context) (<-chan string, <-chan error, error)
}

type SettingsComponent interface {
	SelectSettings(ctx context.Context, title trigger.Title, key storage.Key, value string) error
	SelectedSetting(ctx context.Context, key storage.Key) (string, error)
}
