package transport

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	// EmployeeIDHeader передаёт идентификатор авторизованного аккаунта между transport-слоями.
	EmployeeIDHeader = "x-employee-id"

	// SessionIDHeader передаёт идентификатор авторизованной сессии между transport-слоями.
	SessionIDHeader = "x-session-id"
)

// TokenClaims задаёт общий контракт access token для выпуска и проверки токена.
type TokenClaims struct {
	Identity
	jwt.RegisteredClaims
}

// Identity описывает подтверждённые gateway данные авторизации запроса.
type Identity struct {
	EmployeeID uuid.UUID `json:"employee_id"`
	SessionID  uuid.UUID `json:"session_id"`
}

// IdentityFromContext извлекает данные авторизации из входящей gRPC metadata.
func IdentityFromContext(ctx context.Context) (Identity, error) {
	employeeID, err := metadataUUID(ctx, EmployeeIDHeader)
	if err != nil {
		return Identity{}, fmt.Errorf("employee ID: %w", err)
	}

	sessionID, err := metadataUUID(ctx, SessionIDHeader)
	if err != nil {
		return Identity{}, fmt.Errorf("session ID: %w", err)
	}

	return Identity{
		EmployeeID: employeeID,
		SessionID:  sessionID,
	}, nil
}

func metadataUUID(ctx context.Context, key string) (uuid.UUID, error) {
	value := metadataValue(ctx, key)
	if value == "" {
		return uuid.Nil, fmt.Errorf("not found in incoming metadata")
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse: %w", err)
	}

	return id, nil
}

func metadataValue(ctx context.Context, key string) string {
	values := metadata.ValueFromIncomingContext(ctx, key)
	for i := len(values) - 1; i >= 0; i-- {
		if values[i] != "" {
			return values[i]
		}
	}

	return ""
}
