package services

import (
	"context"
)

func (s *UserService) Logout(ctx context.Context, token string, refreshToken string) (err error) {
	dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().
		RefreshToken_Equal(refreshToken).
		QueryOne(ctx, nil)
	if dbRefreshToken == nil {
		return nil
	}

	err = s.userDB.RefreshToken.Delete(ctx, nil, dbRefreshToken.Id)
	if err != nil {
		return err
	}

	return nil
}
