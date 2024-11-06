package oauth

import (
	"context"
	"fmt"
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
		Client: entity.Client{
			ID:        channel.UserID,
			Secret:    channel.SecretKey,
			Public:    isPublic,
			ClientId:  channel.ClientId,
			XDeviceId: deviceID,
		},
		UserID: channel.UserID,
		TokenInfo: entity.TokenInfo{
			AccessCreateAt: now,
		},
	}

	fmt.Println("now.Add(time.Second * tokenExpiryIn)", now)

	access, refresh, err := u.jwt.Token(context.Background(), data, true)
	token := model.TokenClaim{
		Token:        access,
		RefreshToken: refresh,
	}

	return response.NewSuccessResponse(token, response.StatOK, requestTokenSuccessMessage)
}

func (u *usecase) generateToken(channel *entity.Channel) {

}