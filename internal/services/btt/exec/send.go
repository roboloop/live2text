package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

func (c *client) Exec(ctx context.Context, method string) ([]byte, error) {
	var out bytes.Buffer
	var err error
	cmd := exec.CommandContext(ctx, "osascript", "-e", fmt.Sprintf("tell application \"%s\" to %s", c.bttName, method))
	cmd.Stdout = &out
	defer func() {
		if err != nil {
			c.logger.ErrorContext(ctx, "cannot tell btt", "error", err, "method", method)
		}
	}()
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("command failed: %w", err)
	}

	var unmarshalled []map[string]any
	if err = json.Unmarshal(out.Bytes(), &unmarshalled); err != nil {
		// don't do anything
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var result []map[string]any
	var keysToCheck = []string{"BTTGroupName", "BTTNotes", "BTTGestureNotes"}
	for _, trigger := range unmarshalled {
		for _, key := range keysToCheck {
			v, ok1 := trigger[key]
			_, ok2 := trigger["BTTUUID"].(string)
			if ok1 && ok2 && v == c.appName {
				result = append(result, trigger)
				break
			}
		}
	}

	marshalled, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	c.logger.InfoContext(ctx, "told to the btt", "method", method)

	return marshalled, nil
}
