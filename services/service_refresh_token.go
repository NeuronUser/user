package services

import (
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
	"github.com/NeuronUser/user/models"
	"github.com/NeuronUser/user/storages/user_db"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (s *UserService) RefreshToken(ctx *restful.Context, refreshToken string) (token *models.Token, err error) {
	dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().
		OrderBy(user_db.REFRESH_TOKEN_FIELD_ID, false).
		Limit(0, 1).
		QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbRefreshToken == nil {
		return nil, errors.NotFound("Token已失效，请重新登录")
	}

	if dbRefreshToken.IsLogout == 1 {
		return nil, errors.NotFound("Token已失效，请重新登录")
	}

	if dbRefreshToken.RefreshToken != refreshToken {
		dbRefreshTokenOld, err := s.userDB.RefreshToken.GetQuery().
			RefreshToken_Equal(refreshToken).QueryOne(ctx, nil)
		if err != nil {
			return nil, err
		}

		if dbRefreshTokenOld != nil {
			return nil, errors.NotFound("您已在其它地方登录，请重新登录")
		}

		return nil, errors.NotFound("Token已失效，请重新登录")
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

	dbRefreshToken.RefreshToken = rand.NextHex(16)
	err = s.userDB.RefreshToken.Update(ctx, nil, dbRefreshToken)
	if err != nil {
		return nil, err
	}

	token = &models.Token{}
	token.AccessToken = accessToken
	token.RefreshToken = dbRefreshToken.RefreshToken

	return token, nil
}
