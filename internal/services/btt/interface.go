package btt

import "context"

type Btt interface {
	Initialize(ctx context.Context) error
	Clear(ctx context.Context) error

	SelectedDevice(ctx context.Context) (string, error)
	SelectedLanguage(ctx context.Context) (string, error)
	SelectedFloatingState(ctx context.Context) (string, error)

	SelectDevice(ctx context.Context, device string) error
	SelectLanguage(ctx context.Context, language string) error
	SelectFloatingState(ctx context.Context, floatingState FloatingState) error

	LoadDevices(ctx context.Context) error
	ToggleListening(ctx context.Context) error

	RefreshWidget(ctx context.Context, uuid string) error

	StreamText(ctx context.Context) (<-chan string, <-chan error, error)
	Page() (string, error)
}
