package http

import "context"

type MockClient struct {
	SendResponse [][]byte
	SendError    []error
}

func (m *MockClient) Send(context.Context, string, map[string]any, map[string]string) ([]byte, error) {
	var response []byte
	var err error

	if len(m.SendResponse) > 0 {
		response = m.SendResponse[0]
		m.SendResponse = m.SendResponse[1:]
	}
	if len(m.SendError) > 0 {
		err = m.SendError[0]
		m.SendError = m.SendError[1:]
	}
	return response, err
}
