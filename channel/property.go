package channel

import (
	"time"

	"github.com/sirupsen/logrus"
)

type UsecaseChannelProperty struct {
	ServiceName        string
	Logger             *logrus.Logger
	Location           *time.Location
	ChannelsRepository ChannelsRepository
}
