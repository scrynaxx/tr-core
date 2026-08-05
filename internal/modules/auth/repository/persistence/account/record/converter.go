package record

import (
	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
)

func ToAccount(rec Account) model.Account {
	return model.Account{
		ID:           rec.ID,
		Actor:        rec.Actor,
		Email:        rec.Email,
		PasswordHash: rec.PasswordHash,
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    rec.UpdatedAt,
	}
}
