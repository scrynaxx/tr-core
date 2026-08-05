package account

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
	"github.com/scrynaxx/tr-core/internal/modules/auth/repository"
	"github.com/scrynaxx/tr-core/internal/modules/auth/usecase"
	"github.com/scrynaxx/tr-core/internal/security"
)

type UseCase struct {
	accountRepo repository.Account
}

func New(accountRepo repository.Account) usecase.Account {
	return &UseCase{
		accountRepo: accountRepo,
	}
}

func (uc *UseCase) Create(ctx context.Context, actor security.Actor, email string, password string) (uuid.UUID, error) {
	account, err := model.NewAccount(actor, email, password)
	if err != nil {
		return uuid.Nil(), err
	}

	if err = uc.accountRepo.Create(ctx, account); err != nil {
		return uuid.Nil(), err
	}

	return account.ID, nil
}

func (uc *UseCase) Get(ctx context.Context, accountID uuid.UUID) (model.Account, error) {
	return uc.accountRepo.Get(ctx, accountID)
}
