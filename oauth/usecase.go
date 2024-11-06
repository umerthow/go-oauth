package oauth

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/channel"
	"github.com/umerthow/go-oauth/entity"
	"github.com/umerthow/go-oauth/exception"
	"github.com/umerthow/go-oauth/model"
	"github.com/umerthow/go-oauth/response"
)

const (
	requestTokenSuccessMessage = "Request Token Successfully"
	errorRequestTokenMessage   = "Request Token Failed!"
)

type Usecase interface {
	RequestToken(ctx context.Context, payload model.TokenRequest) response.Response
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

	deviceID := entity.GetDeviceIdFromContext(ctx)
	isPublic := channel.ClientType == "public"

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
			AccessCreateAt: now,
		},
	}

	access, _, err := u.jwt.Token(context.Background(), data, false)
	token := model.TokenClaimResponse{
		Token: access,
	}

	return response.NewSuccessResponse(token, response.StatOK, requestTokenSuccessMessage)
}

func (u *usecase) generateToken(channel *entity.Channel) {

}
