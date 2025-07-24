package background

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrNoSocketFound = errors.New("no socket found")

type SocketManager struct {
	logger    *slog.Logger
	mu        sync.RWMutex
	listeners map[string]net.Listener
	done      map[string]chan struct{}
}

func NewSocketManager(logger *slog.Logger) *SocketManager {
	return &SocketManager{logger: logger, listeners: map[string]net.Listener{}, done: map[string]chan struct{}{}}
}

func (sm *SocketManager) Listen(ctx context.Context, socketPath string, fn func(net.Conn)) error {
	listener, errListen := (&net.ListenConfig{}).Listen(context.Background(), "unix", socketPath)
	if errListen != nil {
		return fmt.Errorf("cannot dial a unix socket: %w", errListen)
	}
	sm.mu.Lock()
	sm.listeners[socketPath] = listener
	sm.done[socketPath] = make(chan struct{})
	sm.mu.Unlock()

	go func() {
		go func() {
			select {
			case <-ctx.Done():
				if err := listener.Close(); err != nil {
					sm.logger.ErrorContext(ctx, "Cannot close listener", "error", err)
				}
			case <-sm.done[socketPath]:
			}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					sm.mu.Lock()
					delete(sm.listeners, socketPath)
					close(sm.done[socketPath])
					sm.mu.Unlock()
					return
				}
				sm.logger.ErrorContext(ctx, "Cannot accept socket connection", "path", socketPath, "error", err)
				continue
			}
			go fn(conn)
		}
	}()

	return nil
}

func (sm *SocketManager) WaitFor(socketPath string) error {
	sm.mu.Lock()
	doneCh, ok := sm.done[socketPath]
	sm.mu.Unlock()
	if !ok {
		return ErrNoSocketFound
	}

	<-doneCh
	return nil
}

func (sm *SocketManager) CloseFor(socketPath string) error {
	sm.mu.Lock()
	listener, ok := sm.listeners[socketPath]
	sm.mu.Unlock()
	if !ok {
		return ErrNoSocketFound
	}

	return listener.Close()
}

func (sm *SocketManager) Close() error {
	g := errgroup.Group{}
	for _, listener := range sm.listeners {
		g.Go(listener.Close)
	}

	return g.Wait()
}

func (sm *SocketManager) TotalOpenSockets() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.listeners)
}
