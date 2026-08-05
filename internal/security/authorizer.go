package security

import (
	"context"
	"time"
	"uuid"

	"github.com/redis/go-redis/v9"
)

type Authorizer interface {
	Authorize(ctx context.Context, token string) (Identity, error)
	RevokeSessions(ctx context.Context, sessionIDs []uuid.UUID, until time.Time) error
}

type authorizer struct {
	redis     redis.UniversalClient
	tokenizer Tokenizer
}

func NewAuthorizer(redis redis.UniversalClient, tokenizer Tokenizer) Authorizer {
	return &authorizer{
		redis:     redis,
		tokenizer: tokenizer,
	}
}

func (a *authorizer) Authorize(ctx context.Context, rawToken string) (Identity, error) {
	if rawToken == "" {
		return Identity{}, ErrInvalidToken
	}

	identity, err := a.tokenizer.ParseAccess(rawToken)
	if err != nil {
		return Identity{}, err
	}

	exists, err := a.redis.Exists(ctx, a.revokedKey(identity.SessionID)).Result()
	if err != nil {
		return Identity{}, err
	}

	if exists > 0 {
		return Identity{}, ErrSessionRevoked
	}

	return identity, nil
}

func (a *authorizer) RevokeSessions(ctx context.Context, sessionIDs []uuid.UUID, until time.Time) error {
	if len(sessionIDs) == 0 {
		return nil
	}

	ttl := until.Sub(time.Now().UTC())
	if ttl <= 0 {
		return nil
	}

	pipe := a.redis.Pipeline()
	for _, id := range sessionIDs {
		pipe.Set(ctx, a.revokedKey(id), "", ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (a *authorizer) revokedKey(sessionID uuid.UUID) string {
	return "access.revoked.session:" + sessionID.String()
}
