package entity

import (
	"time"
)

type GenerateBasic struct {
	ID         string    `json:"id"`
	ClientId   string    `json:"clientId"`
	ClientType string    `json:"clientType"`
	IsActive   bool      `json:"isActive"`
	IsPublic   bool      `json:"isPublic"`
	GrantTypes []string  `json:"grantTypes"`
	Scopes     []string  `json:"scopes"`
	XDeviceId  string    `json:"deviceId"`
	Domain     string    `json:"domain"`
	CreateAt   time.Time `json:"createdAt"`
	TokenInfo  TokenInfo
}

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
