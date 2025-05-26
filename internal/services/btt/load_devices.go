package btt

import (
	"context"
	"fmt"
	"live2text/internal/services/btt/payload"
	"sync"

	"github.com/gordonklaus/portaudio"
)

var mu sync.Mutex

func (b *btt) LoadDevices(ctx context.Context) error {
	var (
		devices         []*portaudio.DeviceInfo
		deviceGroupUUID string
		triggers        []map[string]any
		deviceTriggers  []map[string]any
		err             error
	)

	mu.Lock()
	defer mu.Unlock()

	devices, err = b.audio.List()
	if err != nil {
		return fmt.Errorf("cannot get list of devices: %w", err)
	}

	if deviceGroupUUID, err = b.getTriggerUUID(ctx, deviceGroupTitle); err != nil {
		return fmt.Errorf("cannot find device group: %w", err)
	}

	if triggers, err = b.getTriggers(ctx); err != nil {
		return fmt.Errorf("cannot get triggers: %w", err)
	}
	for _, trigger := range triggers {
		if trigger["BTTTriggerParentUUID"] == deviceGroupUUID &&
			payload.Order(trigger["BTTOrder"].(float64)) > payload.DeviceOrderSelectedDevice {
			deviceTriggers = append(deviceTriggers, trigger)
		}
	}

	for _, trigger := range deviceTriggers {
		if _, err = b.httpClient.Send(ctx, "delete_trigger", nil, map[string]string{"uuid": trigger["BTTUUID"].(string)}); err != nil {
			return fmt.Errorf("cannot delete trigger: %w", err)
		}
	}
	for i, device := range devices {
		var rendered string
		rendered, err = b.renderer.Render(
			"select_device",
			map[string]any{"AppAddress": b.appAddress, "Device": device.Name},
		)
		if err != nil {
			return fmt.Errorf("cannot render select_device script: %w", err)
		}

		devicePayload := make(payload.Payload).
			AddTrigger(device.Name, payload.TriggerTouchBarButton, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
			AddShell(rendered, 0, payload.ShellTypeAdditional)
		if _, err = b.addTrigger(ctx, devicePayload, payload.DeviceOrderSelectedDevice+payload.Order(1+i), deviceGroupUUID); err != nil {
			return fmt.Errorf("cannot create close trigger: %w", err)
		}
	}

	return nil
}
