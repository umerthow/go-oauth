package channel

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/entity"
	"github.com/umerthow/go-oauth/exception"
	"github.com/umerthow/go-oauth/mongodb"
)

type ChannelsRepository interface {
	InsertOne(ctx context.Context, entryData entity.Channel) (err error)
}

type channelRepository struct {
	logger *logrus.Logger
	col    mongodb.Collection
}

func NewChannelRepository(logger *logrus.Logger, db mongodb.Database) ChannelsRepository {
	col := db.Collection("oauth_channel")
	return &channelRepository{logger, col}
}

func (r *channelRepository) InsertOne(ctx context.Context, entryData entity.Channel) (err error) {
	resp, err := r.col.InsertOne(ctx, entryData)
	if err != nil {
		r.logger.Error(err)
		err = exception.ErrInternalServer
		return
	}

	r.logger.Infoln("logId", resp.InsertedID)
	return
}
