package handler

import (
	"go.uber.org/zap"
	"github.com/NeuronFramework/log"
	"github.com/NeuronUser/user/api-private/gen/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

type UserHandler struct {
	logger *zap.Logger
}

func NewUserHandler()(h *UserHandler, err error) {
	h = &UserHandler{}
	h.logger = log.TypedLogger(h)

	return h, nil
}

func (h *UserHandler)GetOauthState(p operations.GetOauthStateParams) (middleware.Responder) {
	return nil
}

func (h *UserHandler)OauthJump(p operations.OauthJumpParams) (middleware.Responder) {
	return nil
}

func (h *UserHandler)Logout(p operations.LogoutParams) (middleware.Responder) {
	return nil
}