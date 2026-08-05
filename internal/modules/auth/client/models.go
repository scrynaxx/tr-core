package authclient

import (
	"time"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/security"
)

type CreateAccountRequest struct {
	Actor    security.Actor
	Email    string
	Password string
}

type Account struct {
	ID        uuid.UUID
	Actor     security.Actor
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
