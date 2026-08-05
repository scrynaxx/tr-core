package event

import (
	"fmt"

	"github.com/scrynaxx/tr-core/contracts/globalevent"
	"github.com/scrynaxx/tr-core/pkg/events"
	"github.com/scrynaxx/tr-core/services/auth/internal/model/event"
	"github.com/scrynaxx/tr-core/services/auth/internal/usecase"
)

func RegisterRoutes(bus *events.Bus, auth usecase.Auth, session usecase.Session) error {
	if err := events.AddSubscriber(bus, "auth", event.SessionRevokedV1, session.HandleRevoked, 1); err != nil {
		return fmt.Errorf("session revoked subscribe: %w", err)
	}

	if err := events.AddSubscriber(bus, "auth", globalevent.EmployeeArchivedV1, auth.HandleEmployeeArchived, 1); err != nil {
		return fmt.Errorf("employee archived subscribe: %w", err)
	}

	if err := events.AddSubscriber(bus, "auth", globalevent.EmployeeCredentialsChangedV1, auth.HandleEmployeeCredentialsChanged, 1); err != nil {
		return fmt.Errorf("employee credentials changed subscribe: %w", err)
	}

	return nil
}
