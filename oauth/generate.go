package oauth

import (
	"context"

	"github.com/umerthow/go-oauth/entity"
)

type (
	// AuthorizeGenerate generate the authorization code interface
	AuthorizeGenerate interface {
		Token(ctx context.Context, data *entity.GenerateBasic) (code string, err error)
	}

	// AccessGenerate generate the access and refresh tokens interface
	AccessGenerate interface {
		Token(ctx context.Context, data *entity.GenerateBasic, isGenRefresh bool) (access, refresh string, err error)
	}
)
