package security

import (
	"errors"
	"uuid"

	"github.com/labstack/echo/v5"
)

type Actor string

const (
	ActorEmployee Actor = "employee"
	ActorCustomer Actor = "customer"
)

type Identity struct {
	AccountID uuid.UUID `json:"account_id"`
	SessionID uuid.UUID `json:"session_id"`
	Actor     Actor     `json:"actor"`
}

const identityKey = "security.identity"

func GetIdentity(c *echo.Context) (Identity, error) {
	identity, ok := c.Get(identityKey).(Identity)
	if !ok {
		return Identity{}, errors.New("identity not found in echo context")
	}

	return identity, nil
}

func SetIdentity(c *echo.Context, identity Identity) {
	c.Set(identityKey, identity)
}
