package entity

import "time"

type TokenInfo struct {
	ClientId         string
	ClientSecret     string
	RedirectURI      string
	AccessExpiresIn  time.Duration
	AccessExpiresAt  time.Time
	AccessCreateAt   time.Time
	RefreshExpiresIn time.Duration
	RefreshExpiresAt time.Time
}

// GetAccessCreateAt create Time
func (t *TokenInfo) GetAccessCreateAt() time.Time {
	return t.AccessCreateAt
}

// GetAccessExpiresIn the lifetime in seconds of the access token
func (t *TokenInfo) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresIn
}
