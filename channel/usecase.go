package channel

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/entity"
	"github.com/umerthow/go-oauth/model"
	"github.com/umerthow/go-oauth/response"
)

const (
	createChannelSuccessMessage = "Create Channel Successfully"
	errorCreateChannelMessage   = "Create Channel Failed!"
	updateChannelSuccessMessage = "Update Channel Successfully"
)

type Usecase interface {
	CreateChannel(ctx context.Context, payload model.RequestChannel) response.Response
	UpdateChannel(ctx context.Context, payload model.RequestChannel, channelID string) response.Response
}

type usecase struct {
	serviceName       string
	logger            *logrus.Logger
	channelRepository ChannelsRepository
	loc               *time.Location
}

func NewChannelUsecase(property UsecaseChannelProperty) *usecase {
	return &usecase{
		serviceName:       property.ServiceName,
		logger:            property.Logger,
		channelRepository: property.ChannelsRepository,
		loc:               property.Location,
	}
}

func (u *usecase) CreateChannel(ctx context.Context, payload model.RequestChannel) response.Response {
	now := time.Now().In(u.loc)

	UserID := uuid.NewString()
	buf := bytes.NewBufferString(UserID)
	buf.WriteString(UserID)

	clientIDGenerate := uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()
	secretKey := base64.URLEncoding.EncodeToString([]byte(clientIDGenerate))

	channel := entity.Channel{
		UserID:      UserID,
		Name:        payload.Name,
		ClientId:    clientIDGenerate,
		IsActive:    true,
		ClientType:  payload.ClientType,
		SecretKey:   strings.ToUpper(strings.TrimRight(secretKey, "=")),
		GrantTypes:  payload.GrantTypes,
		Scopes:      payload.Scopes,
		RedirectURI: payload.RedirectURI,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.channelRepository.InsertOne(ctx, channel); err != nil {
		u.logger.WithContext(ctx).Error(err)

		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, errorCreateChannelMessage)
	}

	return response.NewSuccessResponse("", response.StatOK, createChannelSuccessMessage)
}

func (u *usecase) UpdateChannel(ctx context.Context, payload model.RequestChannel, channelId string) response.Response {

	return response.NewSuccessResponse("", response.StatOK, updateChannelSuccessMessage)
}
