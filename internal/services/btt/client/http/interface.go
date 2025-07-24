package http

import "context"

//go:generate minimock -g -i Client -s _mock.go -o .

type Client interface {
	Send(ctx context.Context, method string, jsonPayload map[string]any, extraPayload map[string]string) ([]byte, error)
}
