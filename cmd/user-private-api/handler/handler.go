package handler

import (
	"context"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/log"
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

func (h *UserHandler) OauthState(p operations.OauthStateParams) middleware.Responder {
	state, err := h.service.OauthState(context.Background(), p.QueryString)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewOauthStateOK().WithPayload(state)
}

func (h *UserHandler) OauthJump(p operations.OauthJumpParams) middleware.Responder {
	result, err := h.service.OauthJump(context.Background(), p.RedirectURI, p.AuthorizationCode, p.State)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewOauthJumpOK().WithPayload(fromOauthJumpResponse(result))
}

func (h *UserHandler) RefreshToken(p operations.RefreshTokenParams) middleware.Responder {
	token, err := h.service.RefreshToken(context.Background(), p.RefreshToken)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewRefreshTokenOK().WithPayload(token)
}

func (h *UserHandler) Logout(p operations.LogoutParams) middleware.Responder {
	p.HTTPRequest.Context()

	err := h.service.Logout(context.Background(), p.Token, p.RefreshToken)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewLogoutOK()
}
