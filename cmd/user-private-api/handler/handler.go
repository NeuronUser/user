package handler

import (
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
	"github.com/NeuronUser/user/api-private/gen/models"
	"github.com/NeuronUser/user/api-private/gen/restapi/operations"
	"github.com/NeuronUser/user/services"
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

func (h *UserHandler) NewOauthState(p operations.NewOauthStateParams) middleware.Responder {
	state, err := h.service.NewOauthState(p.HTTPRequest.Context())
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewNewOauthStateOK().WithPayload(state)
}

func (h *UserHandler) OauthJump(p operations.OauthJumpParams) middleware.Responder {
	token, refreshToken, err := h.service.OauthJump(p.HTTPRequest.Context(), p.AuthorizationCode, p.State)
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewOauthJumpOK().WithPayload(&models.OauthJumpResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *UserHandler) RefreshToken(p operations.RefreshTokenParams) middleware.Responder {
	token, err := h.service.RefreshToken(p.HTTPRequest.Context(), p.RefreshToken)
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewRefreshTokenOK().WithPayload(token)
}

func (h *UserHandler) Logout(p operations.LogoutParams) middleware.Responder {
	p.HTTPRequest.Context()

	err := h.service.Logout(p.HTTPRequest.Context(), p.Token, p.RefreshToken)
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewLogoutOK()
}
