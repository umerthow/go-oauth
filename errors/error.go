package errors

import "errors"

// New returns an error that formats as the given text.
var New = errors.New

// known errors
var (
	ErrInvalidRedirectURI   = errors.New("invalid redirect uri")
	ErrInvalidAuthorizeCode = errors.New("invalid authorize code")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrExpiredAccessToken   = errors.New("expired access token")
	ErrExpiredRefreshToken  = errors.New("expired refresh token")
	ErrMissingCodeVerifier  = errors.New("missing code verifier")
	ErrMissingCodeChallenge = errors.New("missing code challenge")
	ErrInvalidCodeChallenge = errors.New("invalid code challenge")
	ErrUnauthorizedClient   = errors.New("unauthorized_client")
	ErrTokenExpired         = errors.New("token has expired")
	ErrInvalidSignature     = errors.New("token has an invalid signature")
	ErrTokenMalformed       = errors.New("token malformed")
	ErrValidationIssuer     = errors.New("invalid token issuer")
)
