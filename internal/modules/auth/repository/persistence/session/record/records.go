package record

import (
	"time"

	"uuid"
)

type Session struct {
	ID          uuid.UUID `db:"id"`
	AccountID   uuid.UUID `db:"account_id"`
	RefreshHash string    `db:"refresh_hash"`
	UserAgent   string    `db:"user_agent"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
