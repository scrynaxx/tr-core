package record

import "github.com/scrynaxx/tr-core/services/auth/internal/model"

func ToSession(rec Session) model.Session {
	return model.Session{
		ID:          rec.ID,
		EmployeeID:  rec.EmployeeID,
		RefreshHash: rec.RefreshHash,
		UserAgent:   rec.UserAgent,
		ExpiresAt:   rec.ExpiresAt,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}
}
