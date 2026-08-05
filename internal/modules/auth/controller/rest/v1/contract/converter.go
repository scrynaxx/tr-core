package contract

import (
	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
)

func ToAuthResponse(state *model.AuthState) AuthResponse {
	return AuthResponse{
		AccessToken: state.AccessToken,
		ExpiresIn:   state.ExpiresIn,
	}
}
