package audio_wrapper

import (
	"github.com/gordonklaus/portaudio"
	"sync"
)

var initOnce sync.Once
var initError error

func (a *audio) Init() error {
	initOnce.Do(func() {
		initError = portaudio.Initialize()
	})

	return initError
}
