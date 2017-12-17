package services

import (
	"context"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronUser/user/storages/user_db"
)

func (s *UserService) OauthState(ctx context.Context, queryString string) (state string, err error) {
	dbState := &user_db.OauthState{}
	dbState.OauthState = rand.NextHex(16)
	dbState.StateUsed = 0
	dbState.QueryString = queryString
	_, err = s.userDB.OauthState.Insert(ctx, nil, dbState)
	if err != nil {
		return "", err
	}

	return dbState.OauthState, nil
}
