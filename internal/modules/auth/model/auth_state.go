package model

import "time"

type AuthState struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	RefreshUntil time.Time `json:"refresh_until"`
	ExpiresIn    int64     `json:"expires_in"`
}
