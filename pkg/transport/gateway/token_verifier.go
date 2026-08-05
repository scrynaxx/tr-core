package gateway

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/pkg/transport"
)

type tokenVerifier struct {
	parserOptions []jwt.ParserOption
	keyFunc       jwt.Keyfunc
}

func newTokenVerifier(issuer, secret string) *tokenVerifier {
	key := []byte(secret)

	return &tokenVerifier{
		parserOptions: []jwt.ParserOption{
			jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
			jwt.WithIssuer(issuer),
			jwt.WithAudience(issuer),
		},
		keyFunc: func(_ *jwt.Token) (any, error) { return key, nil },
	}
}

func (v *tokenVerifier) Validate(tokenString string) (transport.Identity, error) {
	if tokenString == "" {
		return transport.Identity{}, errors.New("empty token")
	}

	token, err := jwt.ParseWithClaims(tokenString, new(transport.TokenClaims), v.keyFunc, v.parserOptions...)
	if err != nil {
		return transport.Identity{}, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*transport.TokenClaims)
	if !token.Valid || !ok || claims.EmployeeID == uuid.Nil || claims.SessionID == uuid.Nil {
		return transport.Identity{}, errors.New("invalid token claims")
	}

	return claims.Identity, nil
}
