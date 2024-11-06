package oauth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/umerthow/go-oauth/entity"
	err "github.com/umerthow/go-oauth/errors"
)

const (
	issuer = "https://oauth.github.com"
)

// JWTAccessClaims jwt claims
type JWTAccessClaims struct {
	ClientId  string   `json:"clientId"`
	IsActive  bool     `json:"isActive"`
	IsPublic  bool     `json:"isPublic"`
	Scopes    []string `json:"scopes"`
	XDeviceId string   `json:"deviceId"`
	jwt.StandardClaims
}

// NewJWTAccessGenerate create to generate the jwt access token instance
func NewJWTAccessGenerate(kid string, key []byte, method jwt.SigningMethod) *JWTAccessGenerate {
	return &JWTAccessGenerate{
		SignedKeyID:  kid,
		SignedKey:    key,
		SignedMethod: method,
	}
}

// JWTAccessGenerate generate the jwt access token
type JWTAccessGenerate struct {
	SignedKeyID  string
	SignedKey    []byte
	SignedMethod jwt.SigningMethod
}

// Token based on the UUID generated token
func (a *JWTAccessGenerate) Token(ctx context.Context, data *entity.GenerateBasic, isGenRefresh bool) (string, string, error) {
	claims := &JWTAccessClaims{
		ClientId:  data.ClientId,
		Scopes:    data.Scopes,
		IsPublic:  data.IsPublic,
		IsActive:  data.IsActive,
		XDeviceId: data.XDeviceId,
		StandardClaims: jwt.StandardClaims{
			Audience:  data.Domain,
			Issuer:    issuer,
			IssuedAt:  data.TokenInfo.GetAccessCreateAt().Unix(),
			Subject:   data.ID,
			ExpiresAt: data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix(),
		},
	}

	token := jwt.NewWithClaims(a.SignedMethod, claims)
	if a.SignedKeyID != "" {
		token.Header["kid"] = a.SignedKeyID
	}

	access, err := token.SignedString(a.SignedKey)
	if err != nil {
		return "", "", err
	}
	refresh := ""

	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}

func (a *JWTAccessGenerate) Verify(ctx context.Context, accessToken string) (*JWTAccessClaims, error) {
	token, errParse := jwt.ParseWithClaims(accessToken, &JWTAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.SignedKey), nil
	})

	if token.Method != a.SignedMethod {
		return nil, err.ErrInvalidAccessToken
	}

	if errParse != nil {
		if validationErr, ok := errParse.(*jwt.ValidationError); ok {
			if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, err.ErrExpiredAccessToken
			}
			// Check for invalid issuer
			if validationErr.Errors&jwt.ValidationErrorIssuer != 0 {
				return nil, err.ErrValidationIssuer
			}
			// Check for invalid signature
			if validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, err.ErrInvalidSignature
			}
			// Check for malformed token
			if validationErr.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, err.ErrTokenMalformed
			}

			return nil, err.ErrInvalidAccessToken
		} else {
			return nil, errParse
		}
	}

	// Check if the token claims are valid and the token itself is valid
	if claims, ok := token.Claims.(*JWTAccessClaims); ok && token.Valid {
		if claims.Issuer != issuer {
			return nil, err.ErrValidationIssuer
		}

		return claims, nil
	} else {
		return nil, err.ErrInvalidAccessToken
	}
}
