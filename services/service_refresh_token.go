package services

import (
	"context"
	"github.com/NeuronUser/user/storages/user_db"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (tokenString string, err error) {
	dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().RefreshToken_Equal(refreshToken).QueryOne(ctx, nil)
	if err != nil {
		return "", err
	}

	expiresTime := time.Now().Add(time.Hour)

	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   dbRefreshToken.AccountId,
		ExpiresAt: expiresTime.Unix(),
	})
	tokenString, err = userToken.SignedString("0123456789")
	if err != nil {
		return "", err
	}

	dbUserToken := &user_db.UserToken{}
	dbUserToken.AccountId = dbRefreshToken.AccountId
	dbUserToken.UserToken = tokenString
	dbUserToken.ExpiresTime = expiresTime
	_, err = s.userDB.UserToken.Insert(ctx, nil, dbUserToken)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
