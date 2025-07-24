package storage

import "context"

type Key string

//go:generate minimock -g -i Storage -s _mock.go -o .

type Storage interface {
	GetValue(ctx context.Context, key Key) (string, error)
	SetValue(ctx context.Context, key Key, value string) error
}
