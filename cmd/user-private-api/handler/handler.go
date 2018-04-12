package handler

import (
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
	"github.com/NeuronUser/user/api/gen/restapi/operations"
	"github.com/NeuronUser/user/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger  *zap.Logger
	service *services.UserService
}

func NewUserHandler() (h *UserHandler, err error) {
	h = &UserHandler{}
	h.logger = log.TypedLogger(h)
	h.service, err = services.NewUserService()
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *UserHandler) BearerAuth(token string) (userId interface{}, err error) {
	claims := jwt.StandardClaims{}
	_, err = jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("0123456789"), nil
	})
	if err != nil {
		return nil, err
	}

	if claims.Subject == "" {
		return nil, errors.Unknown("验证失败： claims.Subject nil")
	}

	return claims.Subject, nil
}

func (h *UserHandler) GetUserInfo(p operations.GetUserInfoParams, userId interface{}) middleware.Responder {
	userInfo, err := h.service.GetUserInfo(restful.NewContext(p.HTTPRequest), userId.(string))
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewGetUserInfoOK().WithPayload(fromUserInfo(userInfo))
}
