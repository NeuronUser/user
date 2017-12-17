package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful/pointers"
	"github.com/NeuronFramework/sql/wrap"
	"github.com/NeuronFramework2/errors"
	"github.com/NeuronUser/user/models"
	"github.com/NeuronUser/user/remotes/oauth/gen/client/operations"
	"github.com/NeuronUser/user/storages/user_db"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
	"time"
)

func (s *UserService) OauthJump(ctx context.Context, redirectUri string, authorizationCode string, state string) (result *models.OauthJumpResponse, err error) {
	//check state
	dbState, err := s.userDB.OauthState.GetQuery().OauthState_Equal(state).QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbState == nil {
		return nil, errors.InvalidParam("state", "无效的state")
	}

	//get access token
	tokenParams := &operations.TokenParams{}
	tokenParams.Context = ctx
	tokenParams.GrantType = "authorization_code"
	tokenParams.ClientID = pointers.String("100001")
	tokenParams.RedirectURI = pointers.String(redirectUri)
	tokenParams.Code = pointers.String(authorizationCode)
	tokenOk, apiErr := s.oauthClient.Operations.Token(tokenParams, runtime.ClientAuthInfoWriterFunc(
		func(req runtime.ClientRequest, reg strfmt.Registry) error {
			return req.SetHeaderParam("Authorization",
				"Basic "+base64.StdEncoding.EncodeToString(([]byte)("100001"+":"+"100001")))
		}))

	fmt.Println("OauthJump", apiErr)
	if apiErr != nil {
		fmt.Println("OauthJump 111", apiErr)
		return nil, err
	}

	fmt.Println("OauthJump 333")

	accessToken := tokenOk.Payload
	if accessToken == nil {
		return nil, errors.InvalidParam("accessToken", "accessToken nil")
	}

	fmt.Println("OauthJump 444")

	//get account id
	meParams := &operations.MeParams{}
	meParams.Context = ctx
	meParams.AccessToken = accessToken.AccessToken
	meOk, apiError := s.oauthClient.Operations.Me(meParams)
	if apiError != nil {
		fmt.Println("OauthJump ", apiError)
		return nil, apiError
	}
	accountId := meOk.Payload
	if accountId == "" {
		return nil, errors.InvalidParam("accountId", "accountId nil")
	}

	fmt.Println("OauthJump 555")

	expiresTime := time.Now().Add(time.Hour)

	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   accountId,
		ExpiresAt: expiresTime.Unix(),
	})
	tokenString, err := userToken.SignedString([]byte("0123456789"))
	if err != nil {
		return nil, err
	}

	refreshToken := rand.NextHex(16)

	dbOauthTokens := &user_db.OauthTokens{}
	dbOauthTokens.AuthorizationCode = authorizationCode
	dbOauthTokens.AccessToken = accessToken.AccessToken
	dbOauthTokens.RefreshToken = accessToken.RefreshToken
	dbOauthTokens.AccountId = accountId
	_, err = s.userDB.OauthTokens.Insert(ctx, nil, dbOauthTokens)
	if err != nil {
		return nil, err
	}

	dbUserToken := &user_db.UserToken{}
	dbUserToken.AccountId = accountId
	dbUserToken.ExpiresTime = expiresTime
	dbUserToken.UserToken = tokenString
	_, err = s.userDB.UserToken.Insert(ctx, nil, dbUserToken)
	if err != nil {
		return nil, err
	}

	err = s.userDB.TransactionReadCommitted(ctx, func(tx *wrap.Tx) (err error) {
		dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().AccountId_Equal(accountId).QueryOne(ctx, tx)
		if dbRefreshToken == nil {
			dbRefreshToken = &user_db.RefreshToken{}
			dbRefreshToken.AccountId = accountId
			dbRefreshToken.RefreshToken = refreshToken
			_, err = s.userDB.RefreshToken.Insert(ctx, tx, dbRefreshToken)
			if err != nil {
				return err
			}
		} else {
			dbRefreshToken.RefreshToken = refreshToken
			err = s.userDB.RefreshToken.Update(ctx, tx, dbRefreshToken)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	dbState.StateUsed = 1
	err = s.userDB.OauthState.Update(ctx, nil, dbState)
	if err != nil {
		s.logger.Warn("OauthStateUpdate", zap.Error(err))
	}

	result = &models.OauthJumpResponse{}
	result.TokenString = tokenString
	result.RefreshToken = refreshToken
	result.QueryString = dbState.QueryString

	return result, nil
}
