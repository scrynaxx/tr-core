package record

import (
	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
)

func ToSession(rec Session) model.Session {
	return model.Session{
		ID:          rec.ID,
		AccountID:   rec.AccountID,
		RefreshHash: rec.RefreshHash,
		UserAgent:   rec.UserAgent,
		ExpiresAt:   rec.ExpiresAt,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}
}
