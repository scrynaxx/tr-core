package authclient

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/auth/usecase"
)

type Client interface {
	CreateAccount(ctx context.Context, req CreateAccountRequest) (uuid.UUID, error)
	GetAccount(ctx context.Context, accountID uuid.UUID) (Account, error)
}

type client struct {
	accountUc usecase.Account
}

func New(accountUc usecase.Account) Client {
	return &client{
		accountUc: accountUc,
	}
}

func (c *client) CreateAccount(ctx context.Context, req CreateAccountRequest) (uuid.UUID, error) {
	return c.accountUc.Create(ctx, req.Actor, req.Email, req.Password)
}

func (c *client) GetAccount(ctx context.Context, accountID uuid.UUID) (Account, error) {
	account, err := c.accountUc.Get(ctx, accountID)
	if err != nil {
		return Account{}, err
	}

	return Account{
		ID:        account.ID,
		Actor:     account.Actor,
		Email:     account.Email,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}, nil
}
