package text

import "time"

//go:generate minimock -g -i Writer -s _mock.go -o .

type Writer interface {
	PrintFinal(duration time.Duration, text string) error
	PrintCandidate(duration time.Duration, text string) error
	Finalize() error
}
