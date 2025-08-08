package core

import (
	"log/slog"

	"github.com/roboloop/live2text/internal/services"
)

type Server struct {
	logger   *slog.Logger
	services services.Services
}

func NewServer(logger *slog.Logger, services services.Services) *Server {
	return &Server{logger: logger, services: services}
}
