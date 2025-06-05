package btt

import (
	"fmt"
)

func (b *btt) Page() (string, error) {
	page, err := b.renderer.Render("floating_page", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return "", fmt.Errorf("cannot render floating page: %w", err)
	}

	return page, nil
}
