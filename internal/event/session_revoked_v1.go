package event

import (
	"time"

	"uuid"

	"github.com/scrynaxx/tr-core/pkg/messaging"
)

var SessionRevokedV1 = messaging.NewDescriptor[SessionRevokedV1Payload]("auth.session.revoked.v1")

type SessionRevokedV1Payload struct {
	SessionIDs  []uuid.UUID `json:"session_ids"`
	RevokeUntil time.Time   `json:"revoke_until"`
}
