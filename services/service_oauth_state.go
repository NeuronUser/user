package services

import (
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
	"github.com/NeuronUser/user/storages/user_db"
	"time"
)

func (s *UserService) OauthState(ctx *restful.Context, queryString string) (state string, err error) {
	dbRefreshToken := &user_db.RefreshToken{}
	dbRefreshToken.OauthState = rand.NextHex(16)
	dbRefreshToken.QueryString = queryString
	dbRefreshToken.UserAgent = ctx.UserAgent
	dbRefreshToken.GmtLogout = time.Now()
	_, err = s.userDB.RefreshToken.Insert(ctx, nil, dbRefreshToken)
	if err != nil {
		return "", err
	}

	return dbRefreshToken.OauthState, nil
}
