package channel

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/model"
	"github.com/umerthow/go-oauth/response"
)

const (
	createChannelSuccessMessage = "Create Channel Successfully"
	updateChannelSuccessMessage = "Update Channel Successfully"
)

type Usecase interface {
	CreateChannel(ctx context.Context, payload model.Channel) response.Response
	UpdateChannel(ctx context.Context, payload model.Channel, channelID string) response.Response
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

func (u *usecase) CreateChannel(ctx context.Context, payload model.Channel) response.Response {

	return response.NewSuccessResponse("", response.StatOK, createChannelSuccessMessage)
}

func (u *usecase) UpdateChannel(ctx context.Context, payload model.Channel, channelId string) response.Response {

	return response.NewSuccessResponse("", response.StatOK, updateChannelSuccessMessage)
}
