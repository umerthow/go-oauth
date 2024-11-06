package oauth

import (
	"context"
	"net/http"
	"time"

	"github.com/umerthow/go-oauth/entity"
	"github.com/umerthow/go-oauth/model"
)

type (
	// GenerateBasic provide the basis of the generated token data
	GenerateBasic struct {
		Client    entity.Client
		UserID    string
		CreateAt  time.Time
		TokenInfo model.TokenClaim
		Request   *http.Request
	}

	// AuthorizeGenerate generate the authorization code interface
	AuthorizeGenerate interface {
		Token(ctx context.Context, data *GenerateBasic) (code string, err error)
	}

	// AccessGenerate generate the access and refresh tokens interface
	AccessGenerate interface {
		Token(ctx context.Context, data *GenerateBasic, isGenRefresh bool) (access, refresh string, err error)
	}
)
