package channel

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/entity"
	"github.com/umerthow/go-oauth/exception"
	"github.com/umerthow/go-oauth/model"
	"github.com/umerthow/go-oauth/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChannelsRepository interface {
	InsertOne(ctx context.Context, entryData entity.Channel) (err error)
	FindOne(ctx context.Context, payload model.TokenRequest) (channel entity.Channel, err error)
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

func (r *channelRepository) FindOne(ctx context.Context, payload model.TokenRequest) (channel entity.Channel, err error) {
	filter := bson.M{
		"client_id": payload.ClientId,
	}

	if err = r.col.FindOne(ctx, filter).Decode(&channel); err != nil {
		if err != mongo.ErrNoDocuments {
			r.logger.Error(err)
			err = exception.ErrInternalServer
			return
		}
		err = exception.ErrNotFound
		return
	}

	return
}
