package btt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"live2text/internal/services/btt/payload"
)

const (
	bttWidgetName  = "BTTWidgetName"
	bttDescription = "BTTTriggerTypeDescription"
	bttUUID        = "BTTUUID"
)

func (b *btt) getTriggers(ctx context.Context, parentUUID string) ([]map[string]any, error) {
	extraPayload := make(map[string]string)
	if parentUUID != "" {
		extraPayload["trigger_parent_uuid"] = parentUUID
	}

	raw, err := b.httpClient.Send(ctx, "get_triggers", nil, extraPayload)
	if err != nil {
		return nil, fmt.Errorf("cannot query get_triggers: %w", err)
	}

	var unmarshalled []map[string]any
	if err = json.Unmarshal(raw, &unmarshalled); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var result []map[string]any
	var keysToCheck = []string{"BTTGroupName", "BTTNotes", "BTTGestureNotes"}
	for _, trigger := range unmarshalled {
		for _, key := range keysToCheck {
			v, ok1 := trigger[key]
			_, ok2 := trigger[bttUUID].(string)
			if ok1 && ok2 && v == payload.LabelName {
				result = append(result, trigger)
				break
			}
		}
	}

	return result, nil
}

func (b *btt) getTrigger(ctx context.Context, name string) (map[string]any, error) {
	triggers, err := b.getTriggers(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("cannot get triggers: %w", err)
	}

	for _, trigger := range triggers {
		if val, ok1 := trigger[bttDescription]; ok1 && val == name {
			return trigger, nil
		}
	}

	return nil, errors.New("cannot find trigger")
}

func (b *btt) getTriggerUUID(ctx context.Context, name string) (string, error) {
	trigger, err := b.getTrigger(ctx, name)
	if err != nil {
		return "", err
	}

	return trigger[bttUUID].(string), nil
}
