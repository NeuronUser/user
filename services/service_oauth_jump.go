package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/sql/wrap"
	"github.com/NeuronUser/user/remotes/oauth/gen/client/operations"
	"github.com/NeuronUser/user/storages/user_db"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"time"
)

func (s *UserService) OauthJump(ctx context.Context, authorizationCode string, state string) (tokenString string, refreshToken string, err error) {
	tokenOk, err := s.oauthClient.Operations.Token(&operations.TokenParams{}, runtime.ClientAuthInfoWriterFunc(
		func(req runtime.ClientRequest, reg strfmt.Registry) error {
			return req.SetHeaderParam("Authorization",
				"Basic "+base64.StdEncoding.EncodeToString(([]byte)("name"+":"+"password")))
		}))
	if err != nil {
		return "", "", err
	}

	accessToken := tokenOk.Payload
	if accessToken == nil {
		return "", "", fmt.Errorf("accessToken nil")
	}

	meOk, err := s.oauthClient.Operations.Me(&operations.MeParams{})
	if err != nil {
		return "", "", err
	}

	accountId := meOk.Payload
	if accountId == "" {
		return "", "", fmt.Errorf("accountId nil")
	}

	expiresTime := time.Now().Add(time.Hour)

	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   accountId,
		ExpiresAt: expiresTime.Unix(),
	})
	tokenString, err = userToken.SignedString("0123456789")
	if err != nil {
		return "", "", err
	}

	refreshToken = rand.NextBase64(16)

	dbOauthTokens := &user_db.OauthTokens{}
	dbOauthTokens.AuthorizationCode = authorizationCode
	dbOauthTokens.AccessToken = accessToken.AccessToken
	dbOauthTokens.RefreshToken = accessToken.RefreshToken
	dbOauthTokens.AccountId = accountId
	_, err = s.userDB.OauthTokens.Insert(ctx, nil, dbOauthTokens)
	if err != nil {
		return "", "", err
	}

	dbUserToken := &user_db.UserToken{}
	dbUserToken.AccountId = accountId
	dbUserToken.ExpiresTime = expiresTime
	dbUserToken.UserToken = tokenString
	_, err = s.userDB.UserToken.Insert(ctx, nil, dbUserToken)
	if err != nil {
		return "", "", err
	}

	err = s.userDB.TransactionReadCommitted(ctx, func(tx *wrap.Tx) (err error) {
		dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().AccountId_Equal(accountId).QueryOne(ctx, tx)
		if dbRefreshToken == nil {
			dbRefreshToken = &user_db.RefreshToken{}
			dbRefreshToken.RefreshToken = refreshToken
			_, err = s.userDB.RefreshToken.Insert(ctx, tx, dbRefreshToken)
			if err != nil {
				return err
			}
		} else {
			dbRefreshToken.AccountId = accountId
			dbRefreshToken.RefreshToken = refreshToken
			err = s.userDB.RefreshToken.Update(ctx, tx, dbRefreshToken)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return "", "", err
	}

	return tokenString, "", nil
}
