package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessLifetime  = 15 * time.Minute
	RefreshLifetime = 30 * 24 * time.Hour
)

type Tokenizer interface {
	CreateAccess(accountID, sessionID uuid.UUID, actor Actor) (string, error)
	ParseAccess(rawToken string) (Identity, error)
	CreateRefresh() (string, error)
	HashRefresh(refresh string) string
}

type tokenizer struct {
	issuer        string
	secret        []byte
	parserOptions []jwt.ParserOption
}

func NewTokenizer(issuer, secret string) Tokenizer {
	return &tokenizer{
		issuer: issuer,
		secret: []byte(secret),
		parserOptions: []jwt.ParserOption{
			jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
			jwt.WithIssuer(issuer),
		},
	}
}

func (t *tokenizer) CreateAccess(accountID, sessionID uuid.UUID, actor Actor) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		Identity: Identity{
			AccountID: accountID,
			SessionID: sessionID,
			Actor:     actor,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessLifetime)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    t.issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(t.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (t *tokenizer) ParseAccess(rawToken string) (Identity, error) {
	token, err := jwt.ParseWithClaims(rawToken, new(Claims), func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS512 {
			return nil, ErrInvalidToken
		}

		return t.secret, nil
	}, t.parserOptions...)
	if err != nil || !token.Valid {
		return Identity{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return Identity{}, ErrInvalidToken
	}

	return claims.Identity, nil
}

func (t *tokenizer) CreateRefresh() (string, error) {
	value := make([]byte, 64)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(value), nil
}

func (t *tokenizer) HashRefresh(refresh string) string {
	hash := hmac.New(sha256.New, t.secret)
	hash.Write([]byte(refresh))

	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}
