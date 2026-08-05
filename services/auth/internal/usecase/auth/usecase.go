package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/contracts/generated/services/inner/employee"
	"github.com/scrynaxx/tr-core/contracts/globalevent"
	"github.com/scrynaxx/tr-core/pkg/database"
	"github.com/scrynaxx/tr-core/pkg/events"
	"github.com/scrynaxx/tr-core/pkg/transport"
	"github.com/scrynaxx/tr-core/services/auth/internal/model"
	"github.com/scrynaxx/tr-core/services/auth/internal/model/event"
	"github.com/scrynaxx/tr-core/services/auth/internal/repository"
	"github.com/scrynaxx/tr-core/services/auth/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

const (
	// AccessTokenLifetime время жизни access token.
	AccessTokenLifetime = 15 * time.Minute

	// RefreshTokenLifetime время жизни применяется к refresh token и серверной сессии.
	RefreshTokenLifetime = 15 * 24 * time.Hour
)

type UseCase struct {
	issuer         string
	secret         string
	sessionRepo    repository.Session
	employeeClient employee.EmployeeServiceClient
	eventRepo      events.OutboxRepository
	transactor     database.Transactor
}

func NewUseCase(
	issuer string,
	secret string,
	sessionRepo repository.Session,
	employeeClient employee.EmployeeServiceClient,
	eventRepo events.OutboxRepository,
	transactor database.Transactor,
) usecase.Auth {
	return withTracing(&UseCase{
		issuer:         issuer,
		secret:         secret,
		sessionRepo:    sessionRepo,
		employeeClient: employeeClient,
		eventRepo:      eventRepo,
		transactor:     transactor,
	})
}

func (u *UseCase) SignIn(ctx context.Context, email, password, userAgent string) (*model.TokenData, error) {
	credentials, err := u.employeeClient.FindCredentials(ctx, &employee.FindCredentialsRequest{Email: email})
	if err != nil {
		return nil, fmt.Errorf("find credentials: %w", err)
	}

	if credentials == nil || !u.isValidPassword(credentials.PasswordHash, password) {
		return nil, model.ErrInvalidCredentials
	}

	employeeID := uuid.MustParse(credentials.EmployeeId)
	tokens := new(model.TokenData)
	if err = u.transactor.Call(ctx, func(ctx context.Context) error {
		sessionIDs, err := u.sessionRepo.DeleteByUserAgent(ctx, employeeID, userAgent)
		if err != nil {
			return fmt.Errorf("delete session: %w", err)
		}

		for _, id := range sessionIDs {
			if err = u.storeSessionRevoked(ctx, id); err != nil {
				return fmt.Errorf("store session revoked: %w", err)
			}
		}

		refresh, err := u.generateRefresh()
		if err != nil {
			return fmt.Errorf("generate refresh token: %w", err)
		}

		session, err := model.NewSession(employeeID, userAgent, u.hashRefresh(refresh), RefreshTokenLifetime)
		if err != nil {
			return fmt.Errorf("new session: %w", err)
		}

		session, err = u.sessionRepo.Create(ctx, session)
		if err != nil {
			return fmt.Errorf("insert session: %w", err)
		}

		access, err := u.generateAccess(employeeID, session.ID)
		if err != nil {
			return fmt.Errorf("generate access token: %w", err)
		}

		tokens = &model.TokenData{
			AccessToken:  access,
			RefreshToken: refresh,
			ExpiresAt:    session.ExpiresAt,
			ExpiresIn:    int64(AccessTokenLifetime.Seconds()),
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *UseCase) SignOut(ctx context.Context, employeeID, sessionID uuid.UUID) error {
	return u.transactor.Call(ctx, func(ctx context.Context) error {
		affected, err := u.sessionRepo.Delete(ctx, employeeID, sessionID)
		if err != nil {
			return fmt.Errorf("delete session: %w", err)
		}

		if affected == 0 {
			return nil
		}

		if err = u.storeSessionRevoked(ctx, sessionID); err != nil {
			return fmt.Errorf("store session revoked: %w", err)
		}

		return nil
	})
}

func (u *UseCase) Refresh(ctx context.Context, refresh, userAgent string) (*model.TokenData, error) {
	session, err := u.sessionRepo.Get(ctx, u.hashRefresh(refresh), userAgent)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	if session.IsExpired() {
		if err = u.transactor.Call(ctx, func(ctx context.Context) error {
			affected, err := u.sessionRepo.Delete(ctx, session.EmployeeID, session.ID)
			if err != nil {
				return fmt.Errorf("delete session: %w", err)
			}

			if affected > 0 {
				if err = u.storeSessionRevoked(ctx, session.ID); err != nil {
					return fmt.Errorf("store session revoked: %w", err)
				}
			}

			return nil
		}); err != nil {
			return nil, err
		}

		return nil, model.ErrSessionExpired
	}

	refresh, err = u.generateRefresh()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	access, err := u.generateAccess(session.EmployeeID, session.ID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	session.UpdateRefresh(u.hashRefresh(refresh), RefreshTokenLifetime)

	session, err = u.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("update session: %w", err)
	}

	return &model.TokenData{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    session.ExpiresAt,
		ExpiresIn:    int64(AccessTokenLifetime.Seconds()),
	}, nil
}

func (u *UseCase) HandleEmployeeArchived(ctx context.Context, e events.Event[globalevent.EmployeeArchivedDataV1]) error {
	return database.WithTx(ctx, u.transactor, func(ctx context.Context) error {
		sessionIDs, err := u.sessionRepo.DeleteByEmployee(ctx, e.Data.EmployeeID)
		if err != nil {
			return fmt.Errorf("delete employee sessions: %w", err)
		}

		for _, id := range sessionIDs {
			if err = u.storeSessionRevoked(ctx, id); err != nil {
				return fmt.Errorf("store session revoked: %w", err)
			}
		}

		return nil
	})
}

func (u *UseCase) HandleEmployeeCredentialsChanged(ctx context.Context, e events.Event[globalevent.EmployeeCredentialsChangedDataV1]) error {
	return database.WithTx(ctx, u.transactor, func(ctx context.Context) error {
		sessionIDs, err := u.sessionRepo.DeleteByEmployee(ctx, e.Data.EmployeeID)
		if err != nil {
			return fmt.Errorf("delete employee sessions: %w", err)
		}

		for _, id := range sessionIDs {
			if err = u.storeSessionRevoked(ctx, id); err != nil {
				return fmt.Errorf("store session revoked: %w", err)
			}
		}

		return nil
	})
}

func (u *UseCase) generateAccess(employeeID, sessionID uuid.UUID) (string, error) {
	now := time.Now().UTC()

	claims := transport.TokenClaims{
		Identity: transport.Identity{
			EmployeeID: employeeID,
			SessionID:  sessionID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    u.issuer,
			Audience:  []string{u.issuer},
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(u.secret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return token, nil
}

func (u *UseCase) isValidPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (u *UseCase) generateRefresh() (string, error) {
	randomBytes := make([]byte, 64)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func (u *UseCase) hashRefresh(refresh string) string {
	hash := hmac.New(sha256.New, []byte(u.secret))
	_, _ = hash.Write([]byte(refresh))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func (u *UseCase) storeSessionRevoked(ctx context.Context, sessionID uuid.UUID) error {
	e, err := events.NewMessage(event.SessionRevokedV1, event.SessionRevokedDataV1{
		SessionID:   sessionID,
		RevokeUntil: time.Now().UTC().Add(AccessTokenLifetime),
	})
	if err != nil {
		return fmt.Errorf("create session revoked event: %w", err)
	}

	if err = u.eventRepo.StoreEvent(ctx, e); err != nil {
		return fmt.Errorf("store session revoked event: %w", err)
	}

	return nil
}
