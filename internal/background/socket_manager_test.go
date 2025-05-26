package background_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"live2text/internal/background"
	"live2text/internal/utils"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestSocketManager(t *testing.T) {
	ctx := t.Context()
	defaultFn := func(conn net.Conn) {
		conn.Close()
	}

	t.Run("Socket communication between client and server", func(t *testing.T) {
		var (
			ping      = "ping"
			pong      = "pong"
			incoming  string
			outcoming string
			err       error
		)

		sm, socketPath := newSocketManager(ctx)
		defer sm.Close()

		err = sm.Listen(socketPath, func(conn net.Conn) {
			defer conn.Close()
			incoming = read(conn)

			fmt.Fprintf(conn, "%s", pong)
		})
		assertNoError(t, err)

		conn, err := net.Dial("unix", socketPath)
		assertNoError(t, err)
		defer conn.Close()

		_, err = fmt.Fprintf(conn, "%s", ping)
		assertNoError(t, err)

		outcoming = read(conn)
		assertNoError(t, err)

		if incoming != ping {
			t.Errorf("Incoming mismatch: got %v, expected %v", incoming, ping)
		}
		if outcoming != pong {
			t.Errorf("Outcoming mismatch: got %v, expected %v", outcoming, ping)
		}
	})

	t.Run("Listener was closed by socket path", func(t *testing.T) {
		sm, socketPath := newSocketManager(ctx)
		err := sm.Listen(socketPath, defaultFn)
		assertNoError(t, err)

		err = sm.CloseByPath(socketPath)
		assertNoError(t, err)
	})

	t.Run("No listener found to close by socket path", func(t *testing.T) {
		sm, _ := newSocketManager(ctx)

		err := sm.CloseByPath("foo")
		if err != nil && !errors.Is(err, background.ErrNoSocketFound) {
			t.Errorf("ClosedByPath() expected error: got %v, expected %v", err, background.ErrNoSocketFound)
		}
	})

	t.Run("Wait until listener is closed", func(t *testing.T) {
		sm, socketPath := newSocketManager(ctx)
		err := sm.Listen(socketPath, defaultFn)
		assertNoError(t, err)
		sm.Close()
	})

	t.Run("Socket manager status", func(t *testing.T) {
		sm, _ := newSocketManager(ctx)
		sm.Listen("foo", defaultFn)
		sm.Listen("bar", defaultFn)
		totalListeners := sm.Status().TotalListeners
		if totalListeners != 2 {
			t.Errorf("Status() mismatch: got %v, expected %v", totalListeners, 2)
		}

		sm.Close()
	})
}

func newSocketManager(ctx context.Context) (*background.SocketManager, string) {
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d.sock", "live2text", rand.Uint64()))

	return background.NewSocketManager(ctx, utils.NilLogger), path
}

func read(r io.Reader) string {
	buf := make([]byte, 100)
	n, _ := r.Read(buf)
	return string(buf[:n])
}
