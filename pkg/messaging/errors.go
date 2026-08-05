package messaging

import (
	"errors"

	"github.com/scrynaxx/tr-core/pkg/messaging/transport"
)

// ErrLeaseLost сообщает, что outbox worker больше не владеет событием.
var ErrLeaseLost = errors.New("event lease lost")

func isTemporary(err error) bool {
	return errors.Is(err, transport.ErrUnavailable)
}
