package http

import "context"

type Client interface {
	Send(ctx context.Context, method string, jsonPayload map[string]any, extraPayload map[string]string) ([]byte, error)
}
