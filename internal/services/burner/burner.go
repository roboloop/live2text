package burner

import (
	"log/slog"

	"live2text/internal/services/metrics"
)

type burner struct {
	logger  *slog.Logger
	metrics metrics.Metrics
}

func NewBurner(logger *slog.Logger, metrics metrics.Metrics) Burner {
	return &burner{logger: logger, metrics: metrics}
}
