package services

import (
	"context"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronUser/user/models"
	"github.com/NeuronUser/user/storages/user_db"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (token *models.Token, err error) {
	dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().
		RefreshToken_Equal(refreshToken).
		QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}

	if dbRefreshToken == nil {
		return nil, errors.NotFound("refreshToken不存在")
	}

	expiresTime := time.Now().Add(time.Second * models.UserAccessTokenExpireSeconds)

	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   dbRefreshToken.AccountId,
		ExpiresAt: expiresTime.Unix(),
	})
	accessToken, err := userToken.SignedString([]byte("0123456789"))
	if err != nil {
		return nil, err
	}

	dbUserToken := &user_db.UserToken{}
	dbUserToken.AccountId = dbRefreshToken.AccountId
	dbUserToken.UserToken = accessToken
	dbUserToken.ExpiresTime = expiresTime
	_, err = s.userDB.UserToken.Insert(ctx, nil, dbUserToken)
	if err != nil {
		return nil, err
	}

	token = &models.Token{}
	token.AccessToken = accessToken
	token.RefreshToken = dbRefreshToken.RefreshToken

	return token, nil
}
