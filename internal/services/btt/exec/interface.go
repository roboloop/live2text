package exec

import "context"

// Client should be used at all, but it has been placed because webapi doesn't support all methods.
type Client interface {
	Exec(ctx context.Context, method string) ([]byte, error)
}
