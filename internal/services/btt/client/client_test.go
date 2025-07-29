package client_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"live2text/internal/services/btt/client"
	"live2text/internal/services/btt/client/http"
	"live2text/internal/services/btt/client/trigger"
)

func TestGetTriggers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMocks     func(mc *minimock.Controller, c *http.ClientMock)
		parentUUID     trigger.UUID
		label          string
		expectTriggers []trigger.Trigger
		expectErr      string
	}{
		{
			name: "cannot query get_triggers",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			expectErr: "cannot query get_triggers",
		},
		{
			name: "error to parse JSON",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return([]byte("invalid json"), nil)
			},
			expectErr: "error to parse JSON",
		},
		{
			name: "get triggers",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				triggers := `
[
	{
		"BTTGroupName": "app_label",
		"BTTUUID": "trigger1"
	},
	{
		"BTTNotes": "app_label",
		"BTTUUID": "trigger2"
	},
	{
		"BTTGestureNotes": "app_label",
		"BTTUUID": "trigger3"
	},
	{
		"some other key": "app_label",
		"BTTUUID": "trigger4"
	},
	{
		"BTTGroupName": "app_label"
	}
]`
				c.SendMock.Expect(minimock.AnyContext, "get_triggers", nil, map[string]string{
					"trigger_parent_uuid": "1234567890",
				}).Return([]byte(triggers), nil)
			},
			parentUUID: "1234567890",
			label:      "app_label",
			expectTriggers: []trigger.Trigger{
				{
					"BTTGroupName": "app_label",
					"BTTUUID":      "trigger1",
				},
				{
					"BTTNotes": "app_label",
					"BTTUUID":  "trigger2",
				},
				{
					"BTTGestureNotes": "app_label",
					"BTTUUID":         "trigger3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := setupClient(t, tt.setupMocks, tt.label)

			triggers, err := c.GetTriggers(t.Context(), tt.parentUUID)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectTriggers, triggers)
		})
	}
}

func TestGetTrigger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(mc *minimock.Controller, c *http.ClientMock)
		title         trigger.Title
		label         string
		expectTrigger trigger.Trigger
		expectErr     string
	}{
		{
			name: "cannot get triggers",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			expectErr: "cannot get triggers",
		},
		{
			name: "cannot find trigger",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return([]byte("[]"), nil)
			},
			expectErr: "cannot find trigger",
		},
		{
			name: "trigger found",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				triggers := `
[
	{
		"BTTGroupName": "app_label",
		"BTTTriggerTypeDescription": "bar",
		"BTTUUID": "trigger1"
	},
	{
		"BTTNotes": "app_label",
		"BTTTriggerTypeDescription": "foo",
		"BTTUUID": "trigger2"
	}
]`
				c.SendMock.Expect(minimock.AnyContext, "get_triggers", nil, map[string]string{}).
					Return([]byte(triggers), nil)
			},
			title: "foo",
			label: "app_label",
			expectTrigger: map[string]any{
				"BTTNotes":                  "app_label",
				"BTTTriggerTypeDescription": "foo",
				"BTTUUID":                   "trigger2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := setupClient(t, tt.setupMocks, tt.label)

			result, err := c.GetTrigger(t.Context(), tt.title)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectTrigger, result)
		})
	}
}

func TestUpdateTrigger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *http.ClientMock)
		title      trigger.Title
		patch      trigger.Trigger
		label      string
		expectErr  string
	}{
		{
			name: "cannot get trigger",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			title:     "foo",
			patch:     trigger.NewTrigger(),
			expectErr: "cannot get trigger",
		},
		{
			name: "cannot update trigger",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				sendCalls := 0
				c.SendMock.Set(
					func(ctx context.Context, method string, jsonPayload map[string]any, extraPayload map[string]string) ([]byte, error) {
						sendCalls += 1
						switch sendCalls {
						case 1:
							trigger := `
[{
	"BTTNotes":                  "app_label",
	"BTTTriggerTypeDescription": "foo",
	"BTTUUID":                   "trigger2"
}]`
							return []byte(trigger), nil
						default:
							return nil, errors.New("something happened")
						}
					},
				)
			},
			title:     "foo",
			patch:     trigger.NewTrigger(),
			label:     "app_label",
			expectErr: "cannot update trigger",
		},
		{
			name: "trigger updated",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				sendCalls := 0
				c.SendMock.Set(
					func(ctx context.Context, method string, jsonPayload map[string]any, extraPayload map[string]string) ([]byte, error) {
						sendCalls += 1
						switch sendCalls {
						case 1:
							require.Equal(t, "get_triggers", method)
							trigger := `
[{
	"BTTNotes":                  "app_label",
	"BTTTriggerTypeDescription": "foo",
	"BTTUUID":                   "trigger2"
}]`
							return []byte(trigger), nil
						default:
							require.Equal(t, "update_trigger", method)
							return []byte{}, nil
						}
					},
				)
			},
			title: "foo",
			patch: trigger.NewTrigger().AddOrder(42),
			label: "app_label",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := setupClient(t, tt.setupMocks, tt.label)
			err := c.UpdateTrigger(t.Context(), tt.title, tt.patch)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAddTrigger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *http.ClientMock)
		trigger    trigger.Trigger
		label      string
		expectUUID trigger.UUID
		expectErr  string
	}{
		{
			name: "cannot add new trigger",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			trigger:   trigger.NewTrigger(),
			expectErr: "cannot add new trigger",
		},
		{
			name: "trigger added",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Inspect(func(ctx context.Context, method string, jsonPayload map[string]any, extraPayload map[string]string) {
					require.Equal(t, "add_new_trigger", method)
					require.Equal(t, "app_label", jsonPayload["BTTGroupName"])
					require.Equal(t, "foo", extraPayload["parent_uuid"])
				}).
					Return([]byte("trigger1"), nil)
			},
			trigger:    trigger.NewTrigger(),
			label:      "app_label",
			expectUUID: trigger.UUID("trigger1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := setupClient(t, tt.setupMocks, tt.label)
			uuid, err := c.AddTrigger(t.Context(), tt.trigger, "foo")
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectUUID, uuid)
		})
	}
}

func TestDeleteTriggers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *http.ClientMock)
		triggers   []trigger.Trigger
		expectErr  string
	}{
		{
			name: "cannot delete triggers",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			triggers:  []trigger.Trigger{trigger.NewTrigger()},
			expectErr: "cannot delete triggers",
		},
		{
			name: "triggers deleted",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				expectedPayload := map[string]bool{
					"trigger1": false,
					"trigger2": false,
					"trigger3": false,
				}

				c.SendMock.Inspect(func(ctx context.Context, method string, jsonPayload map[string]any, extraPayload map[string]string) {
					require.Equal(t, "delete_trigger", method)
					require.Nil(t, jsonPayload)

					var uuid string
					var ok bool
					if uuid, ok = extraPayload["uuid"]; !ok {
						require.Fail(t, "extra payload does not contain uuid")
					}
					if _, ok = expectedPayload[uuid]; !ok {
						require.Failf(t, "not expected uuid: %s", uuid)
					}
					if expectedPayload[extraPayload["uuid"]] {
						require.Failf(t, "called more than once for uuid: %s", uuid)
					}
					expectedPayload[extraPayload["uuid"]] = true
				}).
					Return([]byte{}, nil)
			},
			triggers: []trigger.Trigger{
				trigger.NewTrigger().AddUUID("trigger1"),
				trigger.NewTrigger().AddUUID("trigger2"),
				trigger.NewTrigger().AddUUID("trigger3"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := setupClient(t, tt.setupMocks, "")

			err := c.DeleteTriggers(t.Context(), tt.triggers)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestRefreshTriggerByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *http.ClientMock)
		title      trigger.Title
		label      string
		expectErr  string
	}{
		{
			name: "cannot get trigger",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			title:     "foo",
			expectErr: "cannot get trigger",
		},
		{
			name: "cannot refresh widget",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				sendCalls := 0
				c.SendMock.Set(
					func(ctx context.Context, method string, jsonPayload map[string]any, extraPayload map[string]string) ([]byte, error) {
						sendCalls += 1
						switch sendCalls {
						case 1:
							trigger := `
[{
	"BTTNotes":                  "app_label",
	"BTTTriggerTypeDescription": "foo",
	"BTTUUID":                   "trigger2"
}]`
							return []byte(trigger), nil
						default:
							return nil, errors.New("something happened")
						}
					},
				)
			},
			title:     "foo",
			label:     "app_label",
			expectErr: "cannot refresh widget",
		},
		{
			name: "trigger refreshed",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				sendCalls := 0
				c.SendMock.Set(
					func(ctx context.Context, method string, jsonPayload map[string]any, extraPayload map[string]string) ([]byte, error) {
						sendCalls += 1
						switch sendCalls {
						case 1:
							trigger := `
[{
	"BTTNotes":                  "app_label",
	"BTTTriggerTypeDescription": "foo",
	"BTTUUID":                   "trigger2"
}]`
							return []byte(trigger), nil
						default:
							return []byte{}, nil
						}
					},
				)
			},
			label: "app_label",
			title: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := setupClient(t, tt.setupMocks, tt.label)

			err := c.RefreshTrigger(t.Context(), tt.title)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestTriggerAction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, c *http.ClientMock)
		action     trigger.Trigger
		expectErr  string
	}{
		{
			name: "cannot trigger action",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			action:    trigger.NewTrigger(),
			expectErr: "cannot trigger action",
		},
		{
			name: "action triggered",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return([]byte(""), nil)
			},
			action:    trigger.NewTrigger(),
			expectErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := setupClient(t, tt.setupMocks, "")
			err := c.TriggerAction(t.Context(), tt.action)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupMocks   func(mc *minimock.Controller, c *http.ClientMock)
		expectHealth bool
	}{
		{
			name: "cannot check health",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("something happened"))
			},
			expectHealth: false,
		},
		{
			name: "health checked successfully",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return(nil, errors.New("unexpected response status code: 400"))
			},
			expectHealth: true,
		},
		{
			name: "health checked unsuccessfully",
			setupMocks: func(mc *minimock.Controller, c *http.ClientMock) {
				c.SendMock.Return([]byte{}, nil)
			},
			expectHealth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := setupClient(t, tt.setupMocks, "")
			health := c.Health(t.Context())
			require.Equal(t, tt.expectHealth, health)
		})
	}
}

func setupClient(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, c *http.ClientMock),
	label string,
) client.Client {
	mc := minimock.NewController(t)
	c := http.NewClientMock(mc)

	if setupMocks != nil {
		setupMocks(mc, c)
	}

	return client.NewClient(c, label)
}
