package exec

import "log/slog"

type client struct {
	logger *slog.Logger

	appName string
	bttName string
}

func NewClient(logger *slog.Logger, appName string, bttName string) Client {
	return &client{logger, appName, bttName}
}
