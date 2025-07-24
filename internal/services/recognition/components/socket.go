package components

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"live2text/internal/background"
	"live2text/internal/services/recognition/text"
)

//go:generate minimock -g -i SocketComponent -s _mock.go -o .

type SocketComponent interface {
	Listen(ctx context.Context, socketPath string, formatter *text.Formatter) error
}

type socketComponent struct {
	logger        *slog.Logger
	socketManager *background.SocketManager
}

func NewSocketComponent(logger *slog.Logger, socketManager *background.SocketManager) SocketComponent {
	return &socketComponent{logger: logger, socketManager: socketManager}
}

func (s *socketComponent) Listen(ctx context.Context, socketPath string, formatter *text.Formatter) error {
	err := s.socketManager.Listen(ctx, socketPath, func(conn net.Conn) {
		defer func() {
			if closeErr := conn.Close(); closeErr != nil {
				s.logger.ErrorContext(ctx, "Error during closing the connection", "error", closeErr)
			}
		}()

		_, _ = conn.Write([]byte(formatter.Format()))
	})
	if err != nil {
		return fmt.Errorf("cannot listen to the socket: %w", err)
	}

	return s.socketManager.WaitFor(socketPath)
}
