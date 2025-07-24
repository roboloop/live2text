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

func NewClient(logger *slog.Logger, bttAddress string, httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: nil,
			Timeout:   time.Second * 10,
		}
	}

	return &client{
		logger: logger,

		client: httpClient,
		bttURL: url.URL{
			Scheme: "http",
			Host:   bttAddress,
		},
	}
}
