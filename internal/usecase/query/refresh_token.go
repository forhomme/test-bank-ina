package query

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/domain/users"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type RefreshTokenHandler decorators.QueryHandler[*users.RefreshToken, *users.Token]

type refreshTokenRepository struct {
	dbRepo users.QueryRepository
}

func NewRefreshTokenRepository(dbRepo users.QueryRepository, logger *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.QueryHandler[*users.RefreshToken, *users.Token] {
	return decorators.ApplyQueryDecorators[*users.RefreshToken, *users.Token](
		refreshTokenRepository{dbRepo: dbRepo},
		logger,
		tracer,
	)
}

func (r refreshTokenRepository) Handle(ctx context.Context, in *users.RefreshToken) (token *users.Token, err error) {
	userData, err := r.dbRepo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	err = in.ParsingToken(userData)
	if err != nil {
		return nil, err
	}
	token, err = userData.GenerateToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}
