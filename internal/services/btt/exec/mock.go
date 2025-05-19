package exec

import "context"

type MockClient struct {
	ExecResponse [][]byte
	ExecError    []error
}

func (m *MockClient) Exec(ctx context.Context, method string) ([]byte, error) {
	var response []byte
	var err error

	if len(m.ExecResponse) > 0 {
		response = m.ExecResponse[0]
		m.ExecResponse = m.ExecResponse[1:]
	}
	if len(m.ExecError) > 0 {
		err = m.ExecError[0]
		m.ExecError = m.ExecError[1:]
	}
	return response, err
}
