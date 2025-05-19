package btt

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	bttWidgetName  = "BTTWidgetName"
	bttDescription = "BTTTriggerTypeDescription"
	bttUuid        = "BTTUUID"
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

	return nil, fmt.Errorf("cannot find trigger")
}

func (b *btt) getTriggerUuid(ctx context.Context, name string) (string, error) {
	trigger, err := b.getTrigger(ctx, name)
	if err != nil {
		return "", err
	}

	return trigger[bttUuid].(string), nil
}
