package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/contracts/globalevent"
	"github.com/scrynaxx/tr-core/pkg/events"
	"github.com/scrynaxx/tr-core/services/auth/internal/model"
	"github.com/scrynaxx/tr-core/services/auth/internal/model/event"
)

type Auth interface {
	// SignIn проверяет реквизиты и создаёт новую сессию.
	SignIn(ctx context.Context, email, password, refresh string) (*model.TokenData, error)

	// SignOut удаляет сессию и сохраняет событие её отзыва.
	SignOut(ctx context.Context, employeeID, sessionID uuid.UUID) error

	// Refresh обновляет refresh token существующей сессии и выпускает новый access token.
	Refresh(ctx context.Context, refresh string, userAgent string) (*model.TokenData, error)

	// HandleEmployeeArchived удаляет все сессии сотрудника и публикует события для отзыва сессий.
	HandleEmployeeArchived(ctx context.Context, e events.Event[globalevent.EmployeeArchivedDataV1]) error

	// HandleEmployeeCredentialsChanged удаляет все сессии сотрудника и публикует события для отзыва сессий.
	HandleEmployeeCredentialsChanged(ctx context.Context, e events.Event[globalevent.EmployeeCredentialsChangedDataV1]) error
}

type Session interface {
	// HandleRevoked сохраняет отозванную сессию в быстром хранилище авторизации.
	HandleRevoked(ctx context.Context, e events.Event[event.SessionRevokedDataV1]) error
}
