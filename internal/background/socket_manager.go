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
	wg        sync.WaitGroup
	listeners map[string]net.Listener
}

type SocketManagerStatus struct {
	TotalListeners int
}

func NewSocketManager(logger *slog.Logger) *SocketManager {
	return &SocketManager{logger: logger, listeners: map[string]net.Listener{}}
}

func (sm *SocketManager) Listen(ctx context.Context, socketPath string, fn func(net.Conn)) error {
	listener, errListen := (&net.ListenConfig{}).Listen(context.Background(), "unix", socketPath)
	if errListen != nil {
		return fmt.Errorf("cannot dial a unix socket: %w", errListen)
	}
	sm.mu.Lock()
	sm.listeners[socketPath] = listener
	sm.mu.Unlock()

	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()

		acceptDoneCh := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				if err := listener.Close(); err != nil {
					sm.logger.ErrorContext(ctx, "Cannot close listener", "error", err)
				}
			case <-acceptDoneCh:
			}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					sm.mu.Lock()
					delete(sm.listeners, socketPath)
					sm.mu.Unlock()
					close(acceptDoneCh)
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

func (sm *SocketManager) CloseByPath(socketPath string) error {
	listener, ok := sm.listeners[socketPath]
	if !ok {
		return ErrNoSocketFound
	}

	if err := listener.Close(); err != nil {
		return err
	}
	return nil
}

func (sm *SocketManager) Close() error {
	g, _ := errgroup.WithContext(context.Background())
	for _, listener := range sm.listeners {
		g.Go(listener.Close)
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (sm *SocketManager) Status() SocketManagerStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return SocketManagerStatus{len(sm.listeners)}
}

func (sm *SocketManager) TotalOpenSockets() float64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return float64(len(sm.listeners))
}
