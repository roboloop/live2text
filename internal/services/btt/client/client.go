package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"live2text/internal/services/btt/client/http"
	"live2text/internal/services/btt/client/trigger"
)

const (
	bttWidgetName  = "BTTWidgetName"
	bttDescription = "BTTTriggerTypeDescription"
	bttUUID        = "BTTUUID"

	bttGroupName    = "BTTGroupName"
	bttNotes        = "BTTNotes"
	bttGestureNotes = "BTTGestureNotes"
)

type client struct {
	httpClient http.Client
	label      string
}

func NewClient(httpClient http.Client, label string) Client {
	return &client{httpClient: httpClient, label: label}
}

func (c *client) GetTriggers(ctx context.Context, parentUUID trigger.UUID) ([]trigger.Trigger, error) {
	extraPayload := make(map[string]string)
	if parentUUID != "" {
		extraPayload["trigger_parent_uuid"] = string(parentUUID)
	}

	raw, err := c.httpClient.Send(ctx, "get_triggers", nil, extraPayload)
	if err != nil {
		return nil, fmt.Errorf("cannot query get_triggers: %w", err)
	}

	var unmarshalled []map[string]any
	if err = json.Unmarshal(raw, &unmarshalled); err != nil {
		return nil, fmt.Errorf("error to parse JSON: %w", err)
	}

	var result []trigger.Trigger
	keysToCheck := []string{bttGroupName, bttNotes, bttGestureNotes}
	for _, t := range unmarshalled {
		for _, key := range keysToCheck {
			v, ok1 := t[key]
			_, ok2 := t[bttUUID].(string)
			if ok1 && ok2 && v == c.label {
				result = append(result, t)
				break
			}
		}
	}

	return result, nil
}

func (c *client) GetTrigger(ctx context.Context, title trigger.Title) (trigger.Trigger, error) {
	triggers, err := c.GetTriggers(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("cannot get triggers: %w", err)
	}

	for _, t := range triggers {
		if val, ok1 := t[bttDescription]; ok1 && val == string(title) {
			return t, nil
		}
	}

	return nil, errors.New("cannot find trigger")
}

func (c *client) UpdateTrigger(ctx context.Context, title trigger.Title, patch trigger.Trigger) error {
	t, err := c.GetTrigger(ctx, title)
	if err != nil {
		return fmt.Errorf("cannot get trigger: %w", err)
	}

	if _, err = c.httpClient.Send(ctx, "update_trigger", patch, map[string]string{"uuid": string(t.UUID())}); err != nil {
		return fmt.Errorf("cannot update trigger: %w", err)
	}

	return nil
}

func (c *client) AddTrigger(ctx context.Context, t trigger.Trigger, parentUUID trigger.UUID) (trigger.UUID, error) {
	extraPayload := make(map[string]string)
	if parentUUID != "" {
		extraPayload["trigger_parent_uuid"] = string(parentUUID)
	}

	uuid, err := c.httpClient.Send(
		ctx,
		"add_new_trigger",
		t.AddLabel(c.label),
		map[string]string{"parent_uuid": string(parentUUID)},
	)
	if err != nil {
		return "", fmt.Errorf("cannot add new trigger: %w", err)
	}

	return trigger.UUID(uuid), nil
}

func (c *client) DeleteTriggers(ctx context.Context, triggers []trigger.Trigger) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(triggers))
	for _, t := range triggers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			jsonPayload := map[string]string{
				"uuid": string(t.UUID()),
			}

			_, err := c.httpClient.Send(ctx, "delete_trigger", nil, jsonPayload)
			errCh <- err
		}()
	}
	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		errs = append(errs, errors.New("cannot delete triggers"))
	}

	return errors.Join(errs...)
}

func (c *client) RefreshTrigger(ctx context.Context, title trigger.Title) error {
	t, err := c.GetTrigger(ctx, title)
	if err != nil {
		return fmt.Errorf("cannot get trigger: %w", err)
	}

	uuid := string(t.UUID())
	if _, err = c.httpClient.Send(ctx, "refresh_widget", nil, map[string]string{"uuid": uuid}); err != nil {
		return fmt.Errorf("cannot refresh widget: %w", err)
	}

	return nil
}

func (c *client) TriggerAction(ctx context.Context, action trigger.Trigger) error {
	if _, err := c.httpClient.Send(ctx, "trigger_action", action, nil); err != nil {
		return fmt.Errorf("cannot trigger action: %w", err)
	}

	return nil
}

func (c *client) Health(ctx context.Context) bool {
	if _, err := c.httpClient.Send(ctx, "health", nil, nil); err != nil {
		return strings.Contains(err.Error(), "unexpected response status code")
	}

	return false
}
