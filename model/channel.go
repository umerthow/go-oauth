package model

import "github.com/umerthow/go-oauth/entity"

type RequestChannel struct {
	Name        string             `json:"name" validate:"required"`
	ClientType  string             `json:"clientType" validate:"oneof=public confidential"`
	GrantTypes  []entity.GrantType `json:"grantTypes" validate:"required"`
	Scopes      []string           `json:"scopes" validate:"required"`
	RedirectURI string             `json:"redirectUri" validate:"required"`
}

type ClientInfo interface {
	GetID() string
	GetSecret() string
	GetDomain() string
	IsPublic() bool
	GetUserID() string
}
