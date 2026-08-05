package event

import (
	"github.com/scrynaxx/tr-core/internal/event"
	"github.com/scrynaxx/tr-core/internal/modules/auth/usecase"
	"github.com/scrynaxx/tr-core/pkg/messaging"
)

func RegisterRoutes(bus *messaging.Bus, authUc usecase.Auth) error {
	if err := messaging.AddSubscriber(bus, "auth", event.EmployeeArchivedV1, authUc.HandleEmployeeArchived, 1); err != nil {
		return err
	}

	if err := messaging.AddSubscriber(bus, "auth", event.SessionRevokedV1, authUc.HandleSessionRevoked, 1); err != nil {
		return err
	}

	return nil
}
