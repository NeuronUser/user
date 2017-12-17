package handler

import (
	api "github.com/NeuronUser/user/api-private/gen/models"
	"github.com/NeuronUser/user/models"
)

func fromOauthJumpResponse(p *models.OauthJumpResponse) (r *api.OauthJumpResponse) {
	if p == nil {
		return nil
	}

	r = &api.OauthJumpResponse{}
	r.Token = p.TokenString
	r.RefreshToken = p.RefreshToken
	r.QueryString = p.QueryString

	return r
}
