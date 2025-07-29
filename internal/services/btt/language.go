package btt

import (
	"context"

	"live2text/internal/services/btt/client/trigger"
	"live2text/internal/services/btt/storage"
)

type languageComponent struct {
	settings SettingsComponent
}

func NewLanguageComponent(settings SettingsComponent) LanguageComponent {
	return &languageComponent{settings: settings}
}

func (b *languageComponent) SelectLanguage(ctx context.Context, language string) error {
	// TODO: restart if it's running?
	return b.settings.SelectSettings(ctx, trigger.TitleSelectedLanguage, storage.SelectedLanguageVariable, language)
}

// SelectedLanguage returns the current selected language. Empty value means missing value.
func (b *languageComponent) SelectedLanguage(ctx context.Context) (string, error) {
	return b.settings.SelectedSetting(ctx, storage.SelectedLanguageVariable)
}
