package services

import (
	"github.com/NeuronFramework/log"
	"github.com/NeuronUser/user/storages/user_db"
	"go.uber.org/zap"
)

type UserService struct {
	logger *zap.Logger
	userDB *user_db.DB
}

func NewUserService() (s *UserService, err error) {
	s = &UserService{}
	s.logger = log.TypedLogger(s)
	s.userDB, err = user_db.NewDB()
	if err != nil {
		return nil, err
	}

	return s, nil
}
