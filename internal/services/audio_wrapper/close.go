package audio_wrapper

import (
	"github.com/gordonklaus/portaudio"
	"sync"
)

var closeOnce sync.Once
var closeError error

func (a *audio) Close() error {
	closeOnce.Do(func() {
		closeError = portaudio.Terminate()
	})

	return closeError
}
