package components_test

import (
	"context"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/background"
	"github.com/roboloop/live2text/internal/services/recognition/components"
	"github.com/roboloop/live2text/internal/services/recognition/text"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestHandleConn(t *testing.T) {
	t.Parallel()

	t.Run("handle conn", func(t *testing.T) {
		t.Parallel()

		socketManager := background.NewSocketManager(logger.NilLogger)
		defer socketManager.Close()
		socketComponent := components.NewSocketComponent(logger.NilLogger, socketManager)

		ctx, cancel := context.WithCancel(t.Context())
		path := filepath.Join(t.TempDir(), "f.sock")
		formatter := text.NewSubtitleFormatter(80, 2)
		formatter.Append("foo", true)

		go func() {
			time.Sleep(20 * time.Millisecond)
			cancel()
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)

			conn, _ := net.Dial("unix", path)
			_ = conn.Close()
		}()

		err := socketComponent.Listen(ctx, path, formatter)

		require.NoError(t, err)
	})

	t.Run("cannot listen to the socket", func(t *testing.T) {
		t.Parallel()

		socketManager := background.NewSocketManager(logger.NilLogger)
		defer socketManager.Close()
		socketComponent := components.NewSocketComponent(logger.NilLogger, socketManager)

		ctx, cancel := context.WithCancel(t.Context())
		formatter := text.NewSubtitleFormatter(80, 2)

		go func() {
			time.Sleep(20 * time.Millisecond)
			cancel()
		}()

		socketPath := filepath.Join(t.TempDir(), strings.Repeat("x", 200))
		err := socketComponent.Listen(ctx, socketPath, formatter)

		require.Error(t, err)
		require.ErrorContains(t, err, "cannot listen to the socket")
	})
}
