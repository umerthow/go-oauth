package model

type TokenRequest struct {
	ClientId     string `json:"clientId"  validate:"required"`
	ClientSecret string `json:"clientSecret" validate:"required"`
	GrantTypes   string `json:"grantTypes" validate:"required"`
}

type TokenClaim struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
