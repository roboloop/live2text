package audiowrapper

import (
	"sync"

	"github.com/gordonklaus/portaudio"
)

var initOnce sync.Once
var errInit error

func (a *audio) Init() error {
	initOnce.Do(func() {
		errInit = portaudio.Initialize()
	})

	return errInit
}
