package storage_test

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/btt/client/http"
	"github.com/roboloop/live2text/internal/services/btt/storage"
)

func TestGetValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMocks  func(mc *minimock.Controller, c *http.ClientMock)
		key         storage.Key
		expectValue string
		expectErr   string
	}{
		{
			name: "cannot get string variable",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			key:       "foo",
			expectErr: "cannot get string variable",
		},
		{
			name: "get string variable",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Expect(minimock.AnyContext, "get_string_variable", nil, map[string]string{
					"variableName": "foo",
				}).Return([]byte("bar"), nil)
			},
			key:         "foo",
			expectValue: "bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := setupStorage(t, tt.setupMocks)

			value, err := s.GetValue(t.Context(), tt.key)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectValue, value)
		})
	}
}

func TestSetValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(*minimock.Controller, *http.ClientMock)
		key        storage.Key
		value      string
		expectErr  string
	}{
		{
			name: "cannot set string variable",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			key:       "foo",
			value:     "bar",
			expectErr: "cannot set string variable",
		},
		{
			name: "set string variable",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Expect(minimock.AnyContext, "set_persistent_string_variable", nil, map[string]string{
					"variableName": "foo",
					"to":           "bar",
				}).Return([]byte{}, nil)
			},
			key:   "foo",
			value: "bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := setupStorage(t, tt.setupMocks)

			err := s.SetValue(t.Context(), tt.key, tt.value)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func setupStorage(t *testing.T, setupMocks func(*minimock.Controller, *http.ClientMock)) storage.Storage {
	mc := minimock.NewController(t)
	c := http.NewClientMock(mc)

	if setupMocks != nil {
		setupMocks(mc, c)
	}

	return storage.NewStorage(c)
}
