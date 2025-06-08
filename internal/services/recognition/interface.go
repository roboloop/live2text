package recognition

import (
	"context"
)

type Recognition interface {
	Start(ctx context.Context, device string, language string) (id string, socketPath string, err error)
	Stop(ctx context.Context, id string) error
	Subs(ctx context.Context, id string) (string, error)
	Has(id string) bool
}
