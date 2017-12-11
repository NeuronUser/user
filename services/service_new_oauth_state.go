package services

import (
	"context"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronUser/user/storages/user_db"
)

func (s *UserService) NewOauthState(ctx context.Context) (state string, err error) {
	dbState := &user_db.OauthState{}
	dbState.OauthState = rand.NextBase64(16)
	dbState.StateUsed = 0
	_, err = s.userDB.OauthState.Insert(ctx, nil, dbState)
	if err != nil {
		return "", err
	}

	return "", nil
}
