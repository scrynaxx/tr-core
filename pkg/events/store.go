package events

import (
	"context"
	"fmt"
)

// Store создаёт сообщение и сохраняет его в outbox с контекстом текущего бизнес-сценария.
func Store[TData any](ctx context.Context, store OutboxRepository, desc Descriptor[TData], data TData) error {
	message, err := NewMessage(desc, data)
	if err != nil {
		return fmt.Errorf("create event message: %w", err)
	}

	if err = store.StoreEvent(ctx, message); err != nil {
		return fmt.Errorf("store event: %w", err)
	}

	return nil
}
