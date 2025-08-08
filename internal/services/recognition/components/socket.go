package components

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/roboloop/live2text/internal/background"
	"github.com/roboloop/live2text/internal/services/recognition/text"
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
	s.logger.InfoContext(ctx, "Listening to the socket", "socketPath", socketPath)

	if err := s.socketManager.Listen(ctx, socketPath, func(conn net.Conn) {
		defer func() {
			if closeErr := conn.Close(); closeErr != nil {
				s.logger.ErrorContext(ctx, "Error during closing the connection", "error", closeErr)
			}
		}()

		// TODO: remove this logic, BTT issue
		formatted := formatter.Format()
		if strings.HasSuffix(formatted, "\n") {
			formatted += " "
		}

		_, _ = conn.Write([]byte(formatted))
	}); err != nil {
		return fmt.Errorf("cannot listen to the socket: %w", err)
	}

	if err := s.socketManager.WaitFor(socketPath); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "Socket listening finished", "socketPath", socketPath)

	return nil
}
