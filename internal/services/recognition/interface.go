package recognition

import (
	"context"
)

//go:generate minimock -g -i Recognition,TaskFactory -s _mock.go -o .

type Recognition interface {
	Start(ctx context.Context, device, language string) (id, socketPath string, err error)
	Stop(ctx context.Context, id string) error
	Text(ctx context.Context, id string) (string, error)
	Has(id string) bool
}

type TaskFactory interface {
	NewTask(device, language string) *Task
}
