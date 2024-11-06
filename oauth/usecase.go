package oauth

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/channel"
	"github.com/umerthow/go-oauth/entity"
	tokenErr "github.com/umerthow/go-oauth/errors"
	"github.com/umerthow/go-oauth/exception"
	"github.com/umerthow/go-oauth/model"
	"github.com/umerthow/go-oauth/response"
)

const (
	requestTokenSuccessMessage       = "Request Token Successfully"
	verifyTokenSuccessMessage        = "Verify Token Successfully"
	errorRequestTokenMessage         = "Request Token Failed!"
	errorNotAllowRequestTokenMessage = "Request Not Allow To Grant Access Token"
)

type Usecase interface {
	RequestToken(ctx context.Context, payload model.TokenRequest) response.Response
	VerifyToken(ctx context.Context, payload model.TokenVerify) response.Response
}

type usecase struct {
	serviceName       string
	logger            *logrus.Logger
	channelRepository channel.ChannelsRepository
	loc               *time.Location
	jwt               JWTAccessGenerate
}

func NewOauthUsecase(property UsecaseOauthProperty) *usecase {
	return &usecase{
		serviceName:       property.ServiceName,
		logger:            property.Logger,
		channelRepository: property.ChannelsRepository,
		loc:               property.Location,
		jwt:               property.JWT,
	}
}

func (u *usecase) RequestToken(ctx context.Context, payload model.TokenRequest) response.Response {
	now := time.Now().In(u.loc)

	channel, err := u.channelRepository.FindOne(ctx, payload)
	if err != nil {
		if err == exception.ErrNotFound {
			return response.NewErrorResponse(exception.ErrForbidden, http.StatusForbidden, nil, response.StatForbidden, err.Error())
		}
		return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, err.Error())
	}

	if channel.SecretKey != payload.ClientSecret {
		return response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, errorRequestTokenMessage)
	}

	if len(channel.GrantTypes) <= 0 {
		return response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, errorNotAllowRequestTokenMessage)
	}

	if channel.GrantTypes[0] != entity.ClientCredentials {
		return response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, errorNotAllowRequestTokenMessage)
	}

	if payload.GrantTypes != entity.ClientCredentials {
		return response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, errorNotAllowRequestTokenMessage)
	}

	deviceID := entity.GetDeviceIdFromContext(ctx)
	isPublic := channel.ClientType == "public"

	tokenExpiryIn := time.Second * 300 // dinamyc expires by client request

	data := &entity.GenerateBasic{
		ID:         channel.ID,
		XDeviceId:  deviceID,
		ClientId:   channel.ClientId,
		ClientType: channel.ClientType,
		IsPublic:   isPublic,
		IsActive:   channel.IsActive,
		GrantTypes: channel.GrantTypes,
		Scopes:     channel.Scopes,
		CreateAt:   channel.CreatedAt,
		Domain:     channel.RedirectURI,
		TokenInfo: entity.TokenInfo{
			AccessCreateAt:  now,
			AccessExpiresIn: tokenExpiryIn,
			AccessExpiresAt: now.Add(tokenExpiryIn),
		},
	}

	access, _, err := u.jwt.Token(context.Background(), data, false)
	token := model.TokenClaimResponse{
		TokenType: "Bearer",
		ExpiredAt: data.TokenInfo.GetAccessExpiresAt(),
		Token:     access,
	}

	return response.NewSuccessResponse(token, response.StatOK, requestTokenSuccessMessage)
}

func (u *usecase) VerifyToken(ctx context.Context, payload model.TokenVerify) response.Response {
	claims, err := u.jwt.Verify(ctx, payload.Token)
	if err != nil {
		if err == tokenErr.ErrExpiredAccessToken {
			return response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatTokenExpired, err.Error())
		}
		return response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, err.Error())
	}

	responseData := model.TokenVerifyResponse{
		ClientId: claims.ClientId,
		Scopes:   claims.Scopes,
	}

	return response.NewSuccessResponse(responseData, response.StatOK, verifyTokenSuccessMessage)
}
