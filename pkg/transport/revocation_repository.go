package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const sessionKeyPrefix = "auth:revoked:session:"

// RevocationRepository хранит отозванные сессии до окончания действия их access token.
type RevocationRepository struct {
	redis redis.Cmdable
}

// NewRevocationRepository создаёт repository отозванных сессий.
func NewRevocationRepository(client redis.Cmdable) *RevocationRepository {
	return &RevocationRepository{
		redis: client,
	}
}

func (r *RevocationRepository) Store(ctx context.Context, sessionID uuid.UUID, ttl time.Duration) error {
	if sessionID == uuid.Nil {
		return fmt.Errorf("session ID is empty")
	}
	if ttl <= 0 {
		return nil
	}

	if err := r.redis.Set(ctx, sessionKey(sessionID), 1, ttl).Err(); err != nil {
		return fmt.Errorf("set revoked session: %w", err)
	}

	return nil
}

func (r *RevocationRepository) IsRevoked(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	if sessionID == uuid.Nil {
		return false, fmt.Errorf("session ID is empty")
	}

	exists, err := r.redis.Exists(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return false, fmt.Errorf("check revoked session: %w", err)
	}

	return exists > 0, nil
}

func sessionKey(sessionID uuid.UUID) string {
	return sessionKeyPrefix + sessionID.String()
}
