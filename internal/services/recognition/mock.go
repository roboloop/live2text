package recognition

import "context"

type MockRecognition struct {
	StartID         string
	StartSocketPath string
	StartError      error

	StopError error

	SubsText  string
	SubsError error
}

func (m *MockRecognition) Start(context.Context, string, string) (string, string, error) {
	return m.StartID, m.StartSocketPath, m.StartError
}

func (m *MockRecognition) Stop(context.Context, string) error {
	return m.StopError
}

func (m *MockRecognition) Subs(context.Context, string) (string, error) {
	return m.SubsText, m.SubsError
}
