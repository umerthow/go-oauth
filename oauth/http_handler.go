package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/umerthow/go-oauth/middleware"
	"github.com/umerthow/go-oauth/model"
	"github.com/umerthow/go-oauth/response"
)

type HTTPHandler struct {
	Logger   *logrus.Logger
	Validate *validator.Validate
	Usecase  Usecase
}

func NewOauthHTTPHandler(logger *logrus.Logger, validate *validator.Validate, router *mux.Router, middleware middleware.RouteMiddleware, usecase Usecase) {
	handler := &HTTPHandler{
		Logger:   logger,
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/go-oauth/v1/token", middleware.Verify(handler.TokenRequest)).Methods(http.MethodPost)
	router.HandleFunc("/go-oauth/v1/token-verification", handler.TokenVerification).Methods(http.MethodGet)
}

func (handler *HTTPHandler) TokenRequest(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var payload model.TokenRequest
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		resp = response.NewErrorResponse(err, http.StatusUnprocessableEntity, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	if err := handler.validateRequestBody(payload); err != nil {
		resp = response.NewErrorResponse(err, http.StatusBadRequest, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	resp = handler.Usecase.RequestToken(ctx, payload)
	response.JSON(w, resp)
}

func (handler *HTTPHandler) TokenVerification(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	queryString := r.URL.Query()
	ctx := r.Context()

	clientId := queryString.Get("clientId")
	token := queryString.Get("token")

	tokenVerify := model.TokenVerify{
		ClientId: clientId,
		Token:    token,
	}
	if tokenVerify.ClientId == "" || tokenVerify.Token == "" {
		err := errors.New("clientId or token can't be empty")
		resp = response.NewErrorResponse(err, http.StatusBadRequest, nil, response.StatusInvalidParameter, err.Error())
		response.JSON(w, resp)
		return
	}

	resp = handler.Usecase.VerifyToken(ctx, tokenVerify)
	response.JSON(w, resp)

}

func (handler *HTTPHandler) validateRequestBody(body interface{}) (err error) {
	err = handler.Validate.Struct(body)
	if err == nil {
		return
	}

	errorFields := err.(validator.ValidationErrors)
	errorField := errorFields[0]
	err = fmt.Errorf("invalid '%s' with value '%v'", errorField.Field(), errorField.Value())

	return
}
