package btt

import (
	"context"
	"fmt"
)

func (b *btt) SelectedLanguage(ctx context.Context) (string, error) {
	value, err := b.getStringVariable(ctx, selectedLanguageVariable)
	if err != nil {
		return "", fmt.Errorf("cannot get selected language variable: %w", err)
	}

	return value, nil
}
