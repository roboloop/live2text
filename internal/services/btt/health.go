package btt

import (
	"context"

	"github.com/roboloop/live2text/internal/services/btt/client"
)

type healthComponent struct {
	client client.Client
}

func NewHealthComponent(client client.Client) HealthComponent {
	return &healthComponent{client: client}
}

func (h *healthComponent) Health(ctx context.Context) bool {
	return h.client.Health(ctx)
}
