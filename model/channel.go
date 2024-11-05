package model

type Channel struct {
	Name        string   `json:"name" validate:"required"`
	ClientType  string   `json:"clientType" validate:"oneof=public confidential"`
	GrantTypes  []string `json:"grantTypes" validate:"required"`
	Scopes      []string `json:"scopes" validate:"required"`
	RedirectURI string   `json:"redirectUri" validate:"required"`
}
