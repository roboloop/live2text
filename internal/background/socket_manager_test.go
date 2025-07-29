package background_test

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"live2text/internal/background"
	"live2text/internal/utils/logger"
)

func TestSocketManager(t *testing.T) {
	t.Parallel()

	defaultFn := func(conn net.Conn) {
		conn.Close()
	}

	t.Run("socket communication between client and server", func(t *testing.T) {
		t.Parallel()

		var (
			ping      = "ping"
			pong      = "pong"
			incoming  string
			outcoming string
			err       error
		)

		sm, socketPath := newSocketManager()
		defer sm.Close()

		err = sm.Listen(t.Context(), socketPath, func(conn net.Conn) {
			defer conn.Close()
			incoming = read(conn)
			fmt.Fprintf(conn, "%s", pong)
		})
		require.NoError(t, err)

		conn, err := net.Dial("unix", socketPath)
		require.NoError(t, err)
		defer conn.Close()

		_, err = fmt.Fprintf(conn, "%s", ping)
		require.NoError(t, err)

		outcoming = read(conn)

		require.Equal(t, ping, incoming, "Incoming mismatch")
		require.Equal(t, pong, outcoming, "Outcoming mismatch")
	})

	t.Run("socket is busy", func(t *testing.T) {
		t.Parallel()

		sm, socketPath := newSocketManager()
		defer sm.Close()
		err := sm.Listen(t.Context(), socketPath, defaultFn)
		require.NoError(t, err)

		err = sm.Listen(t.Context(), socketPath, defaultFn)
		require.Error(t, err)
		require.ErrorIs(t, err, syscall.EADDRINUSE)
	})

	t.Run("listener closed by socket path", func(t *testing.T) {
		t.Parallel()

		sm, socketPath := newSocketManager()
		err := sm.Listen(t.Context(), socketPath, defaultFn)
		require.NoError(t, err)

		err = sm.CloseFor(socketPath)
		require.NoError(t, err)
	})

	t.Run("listener closed by context", func(t *testing.T) {
		t.Parallel()

		sm, socketPath := newSocketManager()
		ctx, cancel := context.WithCancel(t.Context())
		err := sm.Listen(ctx, socketPath, defaultFn)
		require.NoError(t, err)
		cancel()
		time.Sleep(100 * time.Millisecond)
		err = sm.CloseFor(socketPath)
		require.ErrorIs(t, err, background.ErrNoSocketFound)
	})

	t.Run("no listener found to close by socket path", func(t *testing.T) {
		t.Parallel()

		sm, _ := newSocketManager()

		err := sm.CloseFor("foo")
		require.ErrorIs(t, err, background.ErrNoSocketFound)
	})

	t.Run("wait for listener is done", func(t *testing.T) {
		t.Parallel()

		sm, socketPath := newSocketManager()
		ctx, cancel := context.WithCancel(t.Context())
		sm.Listen(ctx, socketPath, defaultFn)
		cancel()

		err := sm.WaitFor(socketPath)

		require.NoError(t, err)
		sm.Close()
	})

	t.Run("unable wait for listener", func(t *testing.T) {
		t.Parallel()

		sm, socketPath := newSocketManager()
		err := sm.WaitFor(socketPath)

		require.Error(t, err)
		require.ErrorIs(t, err, background.ErrNoSocketFound)
		sm.Close()
	})

	t.Run("total open sockets", func(t *testing.T) {
		t.Parallel()

		sm, _ := newSocketManager()
		sm.Listen(t.Context(), "foo", defaultFn)
		sm.Listen(t.Context(), "bar", defaultFn)
		totalOpenSockets := sm.TotalOpenSockets()
		require.Equal(t, 2, totalOpenSockets)

		sm.Close()
	})
}

func newSocketManager() (*background.SocketManager, string) {
	path := filepath.Join(os.TempDir(), fmt.Sprintf("%d.sock", rand.Uint64()))

	return background.NewSocketManager(logger.NilLogger), path
}

func read(r io.Reader) string {
	buf := make([]byte, 100)
	n, _ := r.Read(buf)
	return string(buf[:n])
}
