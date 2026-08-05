package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/samber/lo"
	"github.com/scrynaxx/tr-core/config"
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/auth"
	"github.com/scrynaxx/tr-core/pkg/request"
	"github.com/scrynaxx/tr-core/pkg/transport"
	"github.com/scrynaxx/tr-core/services/auth/internal/model"
	"github.com/scrynaxx/tr-core/services/auth/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	environment config.Environment
	authCase    usecase.Auth
}

func NewController(environment config.Environment, authCase usecase.Auth) pbv1.AuthServiceServer {
	return &Controller{
		environment: environment,
		authCase:    authCase,
	}
}

func (c *Controller) Refresh(ctx context.Context, _ *emptypb.Empty) (*pbv1.AuthResponse, error) {
	refresh := c.getRefreshToken(ctx)
	userAgent := request.HeaderFromContext(ctx, "grpcgateway-user-agent")

	data, err := c.authCase.Refresh(ctx, refresh, userAgent)
	if err != nil {
		return nil, fmt.Errorf("refresh session: %w", err)
	}

	if err = c.applyRefresh(ctx, data); err != nil {
		return nil, fmt.Errorf("set refresh cookie: %w", err)
	}

	return ToAuthResponse(data), nil
}

func (c *Controller) SignIn(ctx context.Context, req *pbv1.SignInRequest) (*pbv1.AuthResponse, error) {
	userAgent := request.HeaderFromContext(ctx, "grpcgateway-user-agent")

	data, err := c.authCase.SignIn(ctx, req.Email, req.Password, userAgent)
	if err != nil {
		return nil, fmt.Errorf("sign in: %w", err)
	}

	if err = c.applyRefresh(ctx, data); err != nil {
		return nil, fmt.Errorf("set refresh cookie: %w", err)
	}

	return ToAuthResponse(data), nil
}

func (c *Controller) SignOut(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	identity, err := transport.IdentityFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get authCase identity: %w", err)
	}

	if err = c.authCase.SignOut(ctx, identity.EmployeeID, identity.SessionID); err != nil {
		return nil, fmt.Errorf("sign out: %w", err)
	}

	if err = c.applyRefresh(ctx, nil); err != nil {
		return nil, fmt.Errorf("clear refresh cookie: %w", err)
	}

	return new(emptypb.Empty), nil
}

func (c *Controller) applyRefresh(ctx context.Context, data *model.TokenData) error {
	md := metadata.New(nil)

	cookie := http.Cookie{
		Name:     "refresh",
		Path:     "/",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: lo.Ternary(c.environment == config.Production, http.SameSiteStrictMode, http.SameSiteNoneMode),
	}
	if data != nil {
		cookie.Value = data.RefreshToken
		cookie.Expires = data.ExpiresAt
		cookie.MaxAge = 0
	}

	md.Append("Set-Cookie", cookie.String())

	if err := grpc.SetHeader(ctx, md); err != nil {
		return fmt.Errorf("set cookie header: %w", err)
	}

	return nil
}

func (c *Controller) getRefreshToken(ctx context.Context) string {
	values := metadata.ValueFromIncomingContext(ctx, "grpcgateway-cookie")
	if len(values) == 0 {
		return ""
	}

	if cookies, err := http.ParseCookie(values[0]); err == nil {
		for i := range cookies {
			if cookies[i].Name == "refresh" {
				return cookies[i].Value
			}
		}
	}

	return ""
}
