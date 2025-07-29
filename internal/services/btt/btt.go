package btt

const (
	defaultInterval = 0.25
)

type btt struct {
	HealthComponent
	InstallingComponent
	ListeningComponent
	DeviceComponent
	LanguageComponent
	ViewModeComponent
	FloatingComponent
	ClipboardComponent
}

func NewBtt(
	healthComponent HealthComponent,
	installingComponent InstallingComponent,
	listeningComponent ListeningComponent,
	deviceComponent DeviceComponent,
	languageComponent LanguageComponent,
	viewModeComponent ViewModeComponent,
	floatingComponent FloatingComponent,
	clipboardComponent ClipboardComponent,
) Btt {
	return &btt{
		HealthComponent:     healthComponent,
		InstallingComponent: installingComponent,
		ListeningComponent:  listeningComponent,
		DeviceComponent:     deviceComponent,
		LanguageComponent:   languageComponent,
		ViewModeComponent:   viewModeComponent,
		FloatingComponent:   floatingComponent,
		ClipboardComponent:  clipboardComponent,
	}
}
