package btt

import (
	"context"
	"fmt"
	"live2text/internal/services/btt/payload"
)

func (b *btt) addTrigger(ctx context.Context, p payload.Payload, order payload.Order, parentUUID string) (string, error) {
	p.AddOrder(order)

	uuid, err := b.httpClient.Send(ctx, "add_new_trigger", p, map[string]string{"parent_uuid": parentUUID})
	if err != nil {
		return "", fmt.Errorf("cannot add new trigger: %w", err)
	}

	return string(uuid), nil
}
