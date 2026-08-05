package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/pkg/events"
)

// SessionRevokedV1 описывает первую версию локального события отзыва сессии.
var SessionRevokedV1 = events.NewDescriptor[SessionRevokedDataV1]("auth.session.revoked.v1")

type SessionRevokedDataV1 struct {
	SessionID uuid.UUID `json:"session_id"`

	RevokeUntil time.Time `json:"revoke_until"`
}
