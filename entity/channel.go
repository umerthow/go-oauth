package entity

import (
	"context"
	"net/http"
	"time"
)

type Channel struct {
	UserID      string    `json:"id" bson:"id"`
	Name        string    `json:"name" bson:"name"`
	ClientId    string    `json:"client_id" bson:"client_id"`
	ClientType  string    `json:"client_type" bson:"client_type"`
	IsActive    bool      `json:"is_active" bson:"is_active"`
	SecretKey   string    `json:"secret_key" bson:"secret_key"`
	GrantTypes  []string  `json:"grant_types" bson:"grant_types"`
	Scopes      []string  `json:"scopes" bson:"scopes"`
	RedirectURI string    `json:"redirect_uri" bson:"redirect_uri"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type Client struct {
	ID        string
	Secret    string
	Domain    string
	Public    bool
	ClientId  string
	XDeviceId string
}

type GenerateBasic struct {
	Client    Client
	UserID    string
	CreateAt  time.Time
	TokenInfo TokenInfo
	Request   *http.Request
}

// GetID client id
func (c *Client) GetID() string {
	return c.ID
}

func GetDeviceIdFromContext(ctx context.Context) string {
	deviceID := ctx.Value(DeviceContextKey{})

	return deviceID.(string)
}
