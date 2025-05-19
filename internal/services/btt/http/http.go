package http

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type client struct {
	logger *slog.Logger
	client *http.Client
	bttURL url.URL
}

func NewClient(logger *slog.Logger, bttAddress string) Client {
	return &client{
		logger: logger,

		client: &http.Client{
			Transport: nil,
			Timeout:   time.Second * 10,
		},
		bttURL: url.URL{
			Scheme: "http",
			Host:   bttAddress,
		},
	}
}
