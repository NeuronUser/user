package handler

import (
	api "github.com/NeuronUser/user/api-private/gen/models"
	"github.com/NeuronUser/user/models"
)

func fromToken(p *models.Token) (r *api.Token) {
	if p == nil {
		return nil
	}

	r = &api.Token{}
	r.AccessToken = &p.AccessToken
	r.RefreshToken = &p.RefreshToken

	return r
}

func fromOauthJumpResponse(p *models.OauthJumpResponse) (r *api.OauthJumpResponse) {
	if p == nil {
		return nil
	}

	r = &api.OauthJumpResponse{}
	r.UserID = &p.UserID
	r.Token = fromToken(p.Token)
	r.QueryString = p.QueryString

	return r
}
