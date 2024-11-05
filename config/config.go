package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func Load() *Config {
	cfg := new(Config)
	cfg.app()
	cfg.logFormatter()
	cfg.mongodb()

	return cfg
}

func (cfg *Config) app() {
	timezone := os.Getenv("TIMEZONE")
	loc, _ := time.LoadLocation("Asia/Jakarta") // default timezone
	appName := os.Getenv("APP_NAME")
	port := os.Getenv("PORT")

	if l, err := time.LoadLocation(timezone); err == nil {
		loc = l
	}
	rawAllowedOrigins := strings.Trim(os.Getenv("ALLOWED_ORIGINS"), " ")

	allowedOrigins := make([]string, 0)
	if rawAllowedOrigins == "" {
		allowedOrigins = append(allowedOrigins, "*")
	} else {
		allowedOrigins = strings.Split(rawAllowedOrigins, ",")
	}

	cfg.Application.Port = port
	cfg.Application.Name = appName
	cfg.Application.AllowedOrigins = allowedOrigins
	cfg.Application.Location = loc
}

func (cfg *Config) logFormatter() {
	formatter := &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			// _, filename := path.Split(f.File)
			filename := fmt.Sprintf("%s:%d", f.File, f.Line)
			return funcname, filename
		},
	}

	cfg.Logger.Formatter = formatter
}

func (cfg *Config) mongodb() {
	appName := os.Getenv("APP_NAME")
	uri := os.Getenv("MONGODB_URL")
	db := os.Getenv("MONGODB_DATABASE")
	minPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MIN_POOL_SIZE"), 10, 64)
	maxPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MAX_POOL_SIZE"), 10, 64)
	maxConnIdleTime, _ := strconv.ParseInt(os.Getenv("MONGODB_MAX_IDLE_CONNECTION_TIME_MS"), 10, 64)
	opts := options.Client().
		ApplyURI(uri).
		SetAppName(appName).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize).
		SetMaxConnIdleTime(time.Millisecond * time.Duration(maxConnIdleTime))

	cfg.Mongodb.ClientOptions = opts
	cfg.Mongodb.Database = db
}
