package btt

import (
	"context"
	"fmt"
)

func (b *btt) RefreshWidget(ctx context.Context, uuid string) error {
	if _, err := b.httpClient.Send(ctx, "refresh_widget", nil, map[string]string{"uuid": uuid}); err != nil {
		return fmt.Errorf("cannot refresh widget: %w", err)
	}

	return nil
}
