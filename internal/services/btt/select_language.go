package btt

import (
	"context"
	"fmt"
)

func (b *btt) SelectLanguage(ctx context.Context, language string) error {
	if err := b.setPersistentStringVariable(ctx, selectedLanguageVariable, language); err != nil {
		return fmt.Errorf("cannot set selected language: %w", err)
	}

	selectedLanguageUUID, err := b.getTriggerUUID(ctx, selectedLanguageTitle)
	if err != nil {
		return fmt.Errorf("cannot get selected language trigger: %w", err)
	}

	if err = b.RefreshWidget(ctx, selectedLanguageUUID); err != nil {
		return fmt.Errorf("cannot refresh selected language widget: %w", err)
	}

	// TODO: restart if it's running?

	return nil
}
