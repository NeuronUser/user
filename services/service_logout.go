package services

import (
	"context"
	"github.com/NeuronFramework/sql/wrap"
)

func (s *UserService) Logout(ctx context.Context, token string, refreshToken string) (err error) {
	err = s.userDB.TransactionReadCommitted(ctx, func(tx *wrap.Tx) (err error) {
		dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().QueryOne(ctx, tx)
		if dbRefreshToken == nil {
			return nil
		}

		err = s.userDB.RefreshToken.Delete(ctx, tx, dbRefreshToken.Id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
