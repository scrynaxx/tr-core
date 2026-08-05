package usecase

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/event"
	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
	"github.com/scrynaxx/tr-core/internal/security"
	"github.com/scrynaxx/tr-core/pkg/messaging"
)

type Account interface {
	Create(ctx context.Context, actor security.Actor, email string, password string) (uuid.UUID, error)
	Get(ctx context.Context, accountID uuid.UUID) (model.Account, error)
}

type Auth interface {
	SignIn(ctx context.Context, actor security.Actor, email, password, userAgent string) (*model.AuthState, error)
	Refresh(ctx context.Context, actor security.Actor, refresh, userAgent string) (*model.AuthState, error)
	SignOut(ctx context.Context, accountID, sessionID uuid.UUID) error
	HandleSessionRevoked(ctx context.Context, e messaging.Event[event.SessionRevokedV1Payload]) error
	HandleEmployeeArchived(ctx context.Context, e messaging.Event[event.EmployeeArchivedV1Payload]) error
}
