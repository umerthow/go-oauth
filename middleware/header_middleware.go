package middleware

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/umerthow/go-oauth/entity"
	"github.com/umerthow/go-oauth/exception"
	"github.com/umerthow/go-oauth/response"
)

const (
	errorForbiddenMessage      = "Forbidden request"
	errorInvalidRequestMessage = "Invalid request"
)

var (
	header   = "Authorization"
	DeviceId = "X-DEVICE-ID"
)

type HeaderValidation struct {
	Logger *logrus.Logger
}

func NewHeaderMiddleware(logger *logrus.Logger) HeaderMiddleware {
	return &HeaderValidation{
		Logger: logger,
	}
}

func (h *HeaderValidation) responseForbidden(w http.ResponseWriter) {
	resp := response.NewErrorResponse(exception.ErrForbidden, http.StatusForbidden, nil, response.StatForbidden, errorForbiddenMessage)
	response.JSON(w, resp)
}

func (h *HeaderValidation) Verify(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceId := r.Header.Get(DeviceId)

		if deviceId == "" {
			h.responseForbidden(w)
			return
		}

		ctx := context.WithValue(r.Context(), entity.DeviceContextKey{}, deviceId)

		r = r.WithContext(ctx)

		next(w, r)

	})

}
