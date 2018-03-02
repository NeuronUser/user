package services

import (
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/log"
	"github.com/NeuronUser/user/remotes/oauth/gen/client"
	"github.com/NeuronUser/user/storages/user_db"
	"go.uber.org/zap"
	"os"
)

type UserService struct {
	logger      *zap.Logger
	userDB      *user_db.DB
	oauthClient *client.Oauth
}

func NewUserService() (s *UserService, err error) {
	s = &UserService{}
	s.logger = log.TypedLogger(s)
	s.userDB, err = user_db.NewDB()
	if err != nil {
		return nil, err
	}

	oauthUrl := os.Getenv("OAUTH_URL")
	if oauthUrl == "" {
		return nil, errors.Unknown("env OAUTH_URL nil")
	}
	s.oauthClient = client.NewHTTPClientWithConfig(nil,
		client.DefaultTransportConfig().WithHost(oauthUrl))

	return s, nil
}
