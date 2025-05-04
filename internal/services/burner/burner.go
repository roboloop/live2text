package burner

import (
	"live2text/internal/services/metrics"
	"log/slog"
)

type burner struct {
	logger  *slog.Logger
	metrics metrics.Metrics
}

func NewBurner(logger *slog.Logger, metrics metrics.Metrics) Burner {
	return &burner{logger, metrics}
}
