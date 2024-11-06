package model

import (
	"time"

	"github.com/umerthow/go-oauth/entity"
)

type TokenRequest struct {
	ClientId     string           `json:"clientId"  validate:"required"`
	ClientSecret string           `json:"clientSecret" validate:"required"`
	GrantTypes   entity.GrantType `json:"grantTypes" validate:"required"`
}

type TokenClaimResponse struct {
	TokenType    string    `json:"tokenType"`
	ExpiredAt    time.Time `json:"expiredAt"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken,omitempty"`
}

type TokenVerify struct {
	ClientId string `json:"clientId"  validate:"required"`
	Token    string `json:"topken" validate:"required"`
}

type TokenVerifyResponse struct {
	ClientId string   `json:"clientId"`
	Scopes   []string `json:"scopes"`
}
