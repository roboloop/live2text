package client

import (
	"context"

	"live2text/internal/services/btt/client/trigger"
)

//go:generate minimock -g -i Client -s _mock.go -o .

type Client interface {
	GetTriggers(ctx context.Context, parentUUID trigger.UUID) ([]trigger.Trigger, error)
	GetTrigger(ctx context.Context, title trigger.Title) (trigger.Trigger, error)
	UpdateTrigger(ctx context.Context, title trigger.Title, patch trigger.Trigger) error
	AddTrigger(ctx context.Context, trigger trigger.Trigger, parentUUID trigger.UUID) (trigger.UUID, error)
	DeleteTriggers(ctx context.Context, triggers []trigger.Trigger) error

	RefreshTrigger(ctx context.Context, title trigger.Title) error
	TriggerAction(ctx context.Context, action trigger.Trigger) error
	Health(ctx context.Context) bool
}
