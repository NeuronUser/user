package services

import (
	"github.com/NeuronUser/user/models"
	"github.com/NeuronUser/user/storages/user_db"
)

func fromUserInfo(p *user_db.User) (r *models.UserInfo) {
	if p == nil {
		return nil
	}

	r = &models.UserInfo{}
	r.UserID = p.UserId
	r.Name = p.UserName
	r.Icon = p.UserIcon

	return r
}
