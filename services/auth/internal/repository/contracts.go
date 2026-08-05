package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/services/auth/internal/model"
)

// Session описывает операции хранения пользовательских сессий.
type Session interface {
	// Get возвращает сессию по хэшу refresh token.
	Get(ctx context.Context, refreshHash, userAgent string) (model.Session, error)

	// Create сохраняет новую сессию.
	Create(ctx context.Context, session model.Session) (model.Session, error)

	// Update обновляет сессию и заменяет предыдущий refresh hash.
	Update(ctx context.Context, session model.Session) (model.Session, error)

	// Delete удаляет конкретную сессию учётной записи.
	Delete(ctx context.Context, employeeID, sessionID uuid.UUID) (int64, error)

	// DeleteByEmployee удаляет все сессии учётной записи.
	DeleteByEmployee(ctx context.Context, employeeID uuid.UUID) ([]uuid.UUID, error)

	// DeleteByUserAgent удаляет сессии, соответствующие user-agent.
	DeleteByUserAgent(ctx context.Context, employeeID uuid.UUID, userAgent string) ([]uuid.UUID, error)
}
