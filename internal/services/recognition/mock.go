package recognition

import "context"

type MockRecognition struct {
	StartId         string
	StartSocketPath string
	StartError      error

	StopError error

	SubsText  string
	SubsError error
}

func (m *MockRecognition) Start(ctx context.Context, device string, language string) (id string, socketPath string, err error) {
	return m.StartId, m.StartSocketPath, m.StartError
}

func (m *MockRecognition) Stop(ctx context.Context, id string) error {
	return m.StopError
}

func (m *MockRecognition) Subs(ctx context.Context, id string) (string, error) {
	return m.SubsText, m.SubsError
}
