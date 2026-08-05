package v1

import (
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/auth"
	"github.com/scrynaxx/tr-core/services/auth/internal/model"
)

func ToAuthResponse(data *model.TokenData) *pbv1.AuthResponse {
	return &pbv1.AuthResponse{
		Token:     data.AccessToken,
		ExpiresIn: data.ExpiresIn,
	}
}
