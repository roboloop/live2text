package btt

const (
	defaultInterval = 0.25
)

type btt struct {
	InitializingComponent
	ListeningComponent
	DeviceComponent
	LanguageComponent
	ViewModeComponent
	FloatingComponent
	ClipboardComponent
}

func NewBtt(
	initializingComponent InitializingComponent,
	listeningComponent ListeningComponent,
	deviceComponent DeviceComponent,
	languageComponent LanguageComponent,
	viewModeComponent ViewModeComponent,
	floatingComponent FloatingComponent,
	clipboardComponent ClipboardComponent,
) Btt {
	return &btt{
		InitializingComponent: initializingComponent,
		ListeningComponent:    listeningComponent,
		DeviceComponent:       deviceComponent,
		LanguageComponent:     languageComponent,
		ViewModeComponent:     viewModeComponent,
		FloatingComponent:     floatingComponent,
		ClipboardComponent:    clipboardComponent,
	}
}
