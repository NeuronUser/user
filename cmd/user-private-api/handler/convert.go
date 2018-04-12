package handler

import (
	api "github.com/NeuronUser/user/api/gen/models"
	"github.com/NeuronUser/user/models"
)

func fromUserInfo(p *models.UserInfo) (r *api.UserInfo) {
	if p == nil {
		return nil
	}

	r = &api.UserInfo{}
	r.UserID = p.UserID
	r.Name = p.Name
	r.Icon = p.Icon

	return r
}
