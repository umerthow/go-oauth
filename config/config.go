package config

import "time"

type Config struct {
	Application struct {
		Port           string
		Name           string
		AllowedOrigins []string
		Location       *time.Location
	}
	Logger struct {
		Formatter logrus.Formatter
	}
	Mongodb struct {
		ClientOptions *options.ClientOptions
		Database      string
	}
}
