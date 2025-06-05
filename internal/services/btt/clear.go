package btt

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

func (b *btt) Clear(ctx context.Context) error {
	var (
		triggers []map[string]any
		uuids    []string
		err      error
	)

	if err = b.triggerCloseGroup(ctx); err != nil {
		return fmt.Errorf("cannot trigger close group: %w", err)
	}

	if triggers, err = b.getTriggers(ctx, ""); err != nil {
		return fmt.Errorf("cannot extract triggers' uuids: %w", err)
	}
	for _, trigger := range triggers {
		uuids = append(uuids, trigger["BTTUUID"].(string))
	}

	if err = b.deleteTriggers(ctx, uuids); err != nil {
		return fmt.Errorf("coudld not delete triggers: %w", err)
	}

	return nil
}

func (b *btt) triggerCloseGroup(ctx context.Context) error {
	payload := map[string]any{
		"BTTPredefinedActionType": 191,
	}

	if _, err := b.httpClient.Send(ctx, "trigger_action", payload, map[string]string{}); err != nil {
		return fmt.Errorf("cannot send json request: %w", err)
	}

	return nil
}

func (b *btt) deleteTriggers(ctx context.Context, uuids []string) error {
	var wg sync.WaitGroup
	var errCh = make(chan error, len(uuids))
	for _, uuid := range uuids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			payload := map[string]string{
				"uuid": uuid,
			}

			_, err := b.httpClient.Send(ctx, "delete_trigger", nil, payload)
			errCh <- err
		}()
	}
	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
