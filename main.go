package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/config"
	"github.com/umerthow/go-oauth/response"
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

	router := mux.NewRouter()
	router.HandleFunc("/go-oauth", index)
}

func index(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSuccessResponse(nil, response.StatOK, indexMessage)
	response.JSON(w, resp)
}
