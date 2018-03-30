package services

import (
	"encoding/base64"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
	"github.com/NeuronUser/user/models"
	"github.com/NeuronUser/user/remotes/oauth/gen/client/operations"
	"github.com/NeuronUser/user/storages/user_db"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"time"
)

func (s *UserService) OauthJump(ctx *restful.Context, redirectUri string, authorizationCode string, state string) (result *models.OauthJumpResponse, err error) {
	dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().OauthState_Equal(state).QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbRefreshToken == nil {
		return nil, errors.InvalidParam("无效的state")
	}

	if dbRefreshToken.AccountId != "" {
		return nil, errors.InvalidParam("state已使用")
	}

	//remote get access token
	tokenParams := &operations.TokenParams{}
	tokenParams.Context = ctx
	tokenParams.GrantType = "authorization_code"
	tokenParams.ClientID = restful.String("100001")
	tokenParams.RedirectURI = restful.String(redirectUri)
	tokenParams.Code = restful.String(authorizationCode)
	tokenOk, apiErr := s.oauthClient.Operations.Token(tokenParams, runtime.ClientAuthInfoWriterFunc(
		func(req runtime.ClientRequest, reg strfmt.Registry) error {
			return req.SetHeaderParam("Authorization",
				"Basic "+base64.StdEncoding.EncodeToString(([]byte)("100001"+":"+"100001")))
		}))
	if apiErr != nil {
		return nil, apiErr
	}
	oauthAccessToken := tokenOk.Payload
	if oauthAccessToken == nil {
		return nil, errors.InvalidParam("获取accessToken失败")
	}

	//remote get account id
	meParams := &operations.MeParams{}
	meParams.Context = ctx
	meParams.AccessToken = oauthAccessToken.AccessToken
	meOk, apiError := s.oauthClient.Operations.Me(meParams)
	if apiError != nil {
		return nil, apiError
	}
	accountId := meOk.Payload
	if accountId == "" {
		return nil, errors.InvalidParam("获取accountId失败")
	}

	//generate user token
	expiresTime := time.Now().Add(time.Second * models.UserAccessTokenExpireSeconds)
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   accountId,
		ExpiresAt: expiresTime.Unix(),
	})
	userAccessToken, err := userToken.SignedString([]byte("0123456789"))
	if err != nil {
		return nil, err
	}
	dbUserToken := &user_db.UserToken{}
	dbUserToken.AccountId = accountId
	dbUserToken.ExpiresTime = expiresTime
	dbUserToken.UserToken = userAccessToken
	_, err = s.userDB.UserToken.Insert(ctx, nil, dbUserToken)
	if err != nil {
		return nil, err
	}

	dbRefreshToken.OauthAuthorizationCode = authorizationCode
	dbRefreshToken.OauthRefreshToken = oauthAccessToken.RefreshToken
	dbRefreshToken.AccountId = accountId
	dbRefreshToken.RefreshToken = rand.NextHex(16)
	s.userDB.RefreshToken.Update(ctx, nil, dbRefreshToken)

	result = &models.OauthJumpResponse{}
	result.UserID = dbUserToken.AccountId
	result.Token = &models.Token{}
	result.Token.AccessToken = userAccessToken
	result.Token.RefreshToken = dbRefreshToken.RefreshToken
	result.QueryString = dbRefreshToken.QueryString

	return result, nil
}
