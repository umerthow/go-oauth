package channel

import (
	"encoding/json"
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

func NewChannelHTTPHandler(logger *logrus.Logger, validate *validator.Validate, router *mux.Router, basicAuth middleware.RouteMiddleware, usecase Usecase) {
	handler := &HTTPHandler{
		Logger:   logger,
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/go-oauth/v1/channel", basicAuth.Verify(handler.CreateChannel)).Methods(http.MethodPost)
}

func (handler *HTTPHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var payload model.RequestChannel
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

	resp = handler.Usecase.CreateChannel(ctx, payload)
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
