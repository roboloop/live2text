package http

import "context"

// TODO: drop one of the methods

type Client interface {
	Send(ctx context.Context, method string, jsonPayload map[string]any, extraPayload map[string]string) ([]byte, error)
}
