package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" // for development
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/channel"
	"github.com/umerthow/go-oauth/config"
	"github.com/umerthow/go-oauth/middleware"
	"github.com/umerthow/go-oauth/mongodb"
	"github.com/umerthow/go-oauth/oauth"
	"github.com/umerthow/go-oauth/response"
	"github.com/umerthow/go-oauth/server"
)

var (
	cfg          *config.Config
	indexMessage string = "Application is running properly"
)

func init() {
	cfg = config.Load()
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)

	// set validator
	vld := validator.New()

	// set mongodb
	mca := mongodb.NewClientAdapter(cfg.Mongodb.ClientOptions)
	if err := mca.Connect(context.Background()); err != nil {
		logger.Fatal(err)
	}

	channelDB := mca.Database(cfg.Mongodb.Database)

	// Basic Auth Initialze Middleware
	// set basic auth middleware
	basicAuthMiddleware := middleware.NewBasicAuth(cfg.BasicAuth.Username, cfg.BasicAuth.Password)
	headerMiddleware := middleware.NewHeaderMiddleware(logger)

	router := mux.NewRouter()
	router.HandleFunc("/go-oauth", index)

	// set cors
	handler := cors.New(cors.Options{
		AllowedOrigins:   cfg.Application.AllowedOrigins,
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	// Channels
	channelRepository := channel.NewChannelRepository(logger, channelDB)
	channelUsecase := channel.NewChannelUsecase(channel.UsecaseChannelProperty{
		ServiceName:        cfg.Application.Name,
		Logger:             logger,
		ChannelsRepository: channelRepository,
		Location:           cfg.Application.Location,
	})

	// Oauth
	oauthUsecase := oauth.NewOauthUsecase(oauth.UsecaseOauthProperty{
		ServiceName:        cfg.Application.Name,
		Logger:             logger,
		ChannelsRepository: channelRepository,
		Location:           cfg.Application.Location,
		JWT: oauth.JWTAccessGenerate{
			SignedKeyID:  uuid.NewString(),
			SignedKey:    []byte(cfg.JWT.PrivateKey),
			SignedMethod: jwt.SigningMethodHS512,
		},
	})

	// Routes Handler
	channel.NewChannelHTTPHandler(logger, vld, router, basicAuthMiddleware, channelUsecase)
	oauth.NewOauthHTTPHandler(logger, vld, router, headerMiddleware, oauthUsecase)

	// initiate server
	srv := server.NewServer(logger, handler, cfg.Application.Port)
	srv.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	// closing service for a gracefull shutdown.
	srv.Close()
	mca.Disconnect(context.Background())

}

func index(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSuccessResponse(nil, response.StatOK, indexMessage)
	response.JSON(w, resp)
}
