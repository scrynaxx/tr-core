package record

import (
	"time"

	"uuid"

	"github.com/scrynaxx/tr-core/internal/security"
)

type Account struct {
	ID           uuid.UUID      `db:"id"`
	Actor        security.Actor `db:"actor"`
	Email        string         `db:"email"`
	PasswordHash string         `db:"password_hash"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}
