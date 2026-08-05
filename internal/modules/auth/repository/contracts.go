package repository

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
	"github.com/scrynaxx/tr-core/internal/security"
)

type Account interface {
	Get(ctx context.Context, accountID uuid.UUID) (model.Account, error)
	GetByEmail(ctx context.Context, actor security.Actor, email string) (model.Account, error)
	Create(ctx context.Context, account model.Account) error
	Update(ctx context.Context, account model.Account) error
}

type Session interface {
	Get(ctx context.Context, actor security.Actor, refreshHash, userAgent string) (model.Session, error)
	Create(ctx context.Context, session model.Session) error
	Update(ctx context.Context, session model.Session) error
	Delete(ctx context.Context, accountID, sessionID uuid.UUID) (bool, error)
	DeleteByAccount(ctx context.Context, accountID uuid.UUID) ([]uuid.UUID, error)
	DeleteByUserAgent(ctx context.Context, accountID uuid.UUID, userAgent string) ([]uuid.UUID, error)
}
