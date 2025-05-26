package btt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	bttWidgetName  = "BTTWidgetName"
	bttDescription = "BTTTriggerTypeDescription"
	bttUUID        = "BTTUUID"
)

func (b *btt) getTrigger(ctx context.Context, name string) (map[string]any, error) {
	out, err := b.execClient.Exec(ctx, "get_triggers")
	if err != nil {
		return nil, fmt.Errorf("cannot exec get_triggers command: %w", err)
	}

	var result []map[string]any
	if err = json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, trigger := range result {
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
