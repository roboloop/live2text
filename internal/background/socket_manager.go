package background

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net"
	"sync"
)

var NoSocketFound = errors.New("no socket found")

type SocketManager struct {
	ctx       context.Context
	logger    *slog.Logger
	mu        sync.Mutex
	wg        sync.WaitGroup
	listeners map[string]net.Listener
}

type SocketManagerStatus struct {
	TotalListeners int
}

func NewSocketManager(ctx context.Context, logger *slog.Logger) *SocketManager {
	return &SocketManager{ctx: ctx, logger: logger, listeners: map[string]net.Listener{}}
}

func (sm *SocketManager) Listen(socketPath string, fn func(net.Conn)) error {
	//socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d.sock", sm.prefix, rand.Uint64()))
	listener, err := (&net.ListenConfig{}).Listen(sm.ctx, "unix", socketPath)
	if err != nil {
		return fmt.Errorf("cannot dial a unix socket: %w", err)
	}
	sm.mu.Lock()
	sm.listeners[socketPath] = listener
	sm.mu.Unlock()

	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					sm.mu.Lock()
					delete(sm.listeners, socketPath)
					sm.mu.Unlock()
					return
				}
				sm.logger.ErrorContext(sm.ctx, "Cannot accept socket connection", "path", socketPath, "error", err)
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
		return NoSocketFound
	}

	if err := listener.Close(); err != nil {
		return err
	}
	return nil
}

func (sm *SocketManager) Close() error {
	g, _ := errgroup.WithContext(context.Background())
	for _, listener := range sm.listeners {
		g.Go(func() error {
			return listener.Close()
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	//sm.wg.Wait()
	return nil
}

func (sm *SocketManager) Status() SocketManagerStatus {
	return SocketManagerStatus{len(sm.listeners)}
}
