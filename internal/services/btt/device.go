package btt

import (
	"context"
	"fmt"
	"sync"

	"live2text/internal/services/audio"
	"live2text/internal/services/btt/client"
	"live2text/internal/services/btt/client/trigger"
	"live2text/internal/services/btt/storage"
	"live2text/internal/services/btt/tmpl"
)

var loadDevicesMutex sync.Mutex

type deviceComponent struct {
	audio    audio.Audio
	client   client.Client
	renderer tmpl.Renderer
	settings SettingsComponent
}

func NewDeviceComponent(
	audio audio.Audio,
	client client.Client,
	renderer tmpl.Renderer,
	settings SettingsComponent,
) DeviceComponent {
	return &deviceComponent{audio: audio, client: client, renderer: renderer, settings: settings}
}

func (d *deviceComponent) LoadDevices(ctx context.Context) error {
	loadDevicesMutex.Lock()
	defer loadDevicesMutex.Unlock()

	deviceDir, err := d.client.GetTrigger(ctx, trigger.TitleDeviceDir)
	if err != nil {
		return fmt.Errorf("cannot get device dir trigger: %w", err)
	}

	if err = d.deleteDevices(ctx, deviceDir.UUID()); err != nil {
		return err
	}

	if err = d.addDevices(ctx, deviceDir.UUID()); err != nil {
		return err
	}

	return nil
}

func (d *deviceComponent) deleteDevices(ctx context.Context, parentUUID trigger.UUID) error {
	triggers, err := d.client.GetTriggers(ctx, parentUUID)
	if err != nil {
		return fmt.Errorf("cannot get device triggers: %w", err)
	}

	var deviceTriggers []trigger.Trigger
	for _, t := range triggers {
		if t.HasTapScript() {
			deviceTriggers = append(deviceTriggers, t)
		}
	}
	if err = d.client.DeleteTriggers(ctx, deviceTriggers); err != nil {
		return fmt.Errorf("cannot delete triggers: %w", err)
	}

	return nil
}

func (d *deviceComponent) addDevices(ctx context.Context, parentUUID trigger.UUID) error {
	devices, err := d.audio.ListOfNames()
	if err != nil {
		return fmt.Errorf("cannot get a list of devices: %w", err)
	}

	selectedDevice, err := d.client.GetTrigger(ctx, trigger.TitleSelectedDevice)
	if err != nil {
		return fmt.Errorf("cannot get selected device trigger: %w", err)
	}

	after := selectedDevice
	for _, device := range devices {
		after = trigger.NewTapButton(trigger.Title(device), d.renderer.SelectDevice(device)).AddOrderAfter(after)
		if _, err = d.client.AddTrigger(ctx, after, parentUUID); err != nil {
			return fmt.Errorf("cannot create trigger: %w", err)
		}
	}

	return nil
}

func (d *deviceComponent) SelectDevice(ctx context.Context, device string) error {
	// TODO: restart if it's running?
	return d.settings.SelectSettings(ctx, trigger.TitleSelectedDevice, storage.SelectedDeviceVariable, device)
}

func (d *deviceComponent) SelectedDevice(ctx context.Context) (string, error) {
	return d.settings.SelectedSetting(ctx, storage.SelectedDeviceVariable)
}
