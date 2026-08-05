package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/pkg/events"
	"github.com/scrynaxx/tr-core/services/auth/internal/model/event"
	"github.com/scrynaxx/tr-core/services/auth/internal/usecase"
)

type RevokeRepository interface {
	Store(ctx context.Context, sessionID uuid.UUID, ttl time.Duration) error
}

type UseCase struct {
	revokeRepo RevokeRepository
}

func NewUseCase(revokeRepo RevokeRepository) usecase.Session {
	return withTracing(&UseCase{
		revokeRepo: revokeRepo,
	})
}

func (u *UseCase) HandleRevoked(ctx context.Context, e events.Event[event.SessionRevokedDataV1]) error {
	ttl := time.Until(e.Data.RevokeUntil)
	if ttl <= 0 {
		return nil
	}

	if err := u.revokeRepo.Store(ctx, e.Data.SessionID, ttl); err != nil {
		return fmt.Errorf("store revoked session: %w", err)
	}

	return nil
}
