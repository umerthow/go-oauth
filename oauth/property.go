package oauth

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/channel"
)

type UsecaseOauthProperty struct {
	ServiceName        string
	Logger             *logrus.Logger
	Location           *time.Location
	ChannelsRepository channel.ChannelsRepository
	JWT                JWTAccessGenerate
}
