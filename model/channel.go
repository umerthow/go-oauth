package model

type RequestChannel struct {
	Name        string   `json:"name" validate:"required"`
	ClientType  string   `json:"clientType" validate:"oneof=public confidential"`
	GrantTypes  []string `json:"grantTypes" validate:"required"`
	Scopes      []string `json:"scopes" validate:"required"`
	RedirectURI string   `json:"redirectUri" validate:"required"`
}

type ClientInfo interface {
	GetID() string
	GetSecret() string
	GetDomain() string
	IsPublic() bool
	GetUserID() string
}
