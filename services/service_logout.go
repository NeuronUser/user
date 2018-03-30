package services

import (
	"github.com/NeuronFramework/restful"
	"time"
)

func (s *UserService) Logout(ctx *restful.Context, token string, refreshToken string) (err error) {
	dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().
		RefreshToken_Equal(refreshToken).
		QueryOne(ctx, nil)
	if dbRefreshToken == nil {
		return nil
	}

	dbRefreshToken.IsLogout = 1
	dbRefreshToken.GmtLogout = time.Now().UTC()
	err = s.userDB.RefreshToken.Update(ctx, nil, dbRefreshToken)
	if err != nil {
		return err
	}

	return nil
}
