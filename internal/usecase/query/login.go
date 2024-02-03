package query

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/common/errs"
	"test_ina_bank/internal/domain/users"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type LoginHandler decorators.QueryHandler[*users.Login, *users.Token]

type loginRepository struct {
	dbRepo users.QueryRepository
}

func NewLoginRepository(dbRepo users.QueryRepository, logger *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.QueryHandler[*users.Login, *users.Token] {
	return decorators.ApplyQueryDecorators[*users.Login, *users.Token](
		loginRepository{dbRepo: dbRepo},
		logger,
		tracer,
	)
}

func (l loginRepository) Handle(ctx context.Context, in *users.Login) (token *users.Token, err error) {
	userData, err := l.dbRepo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	if !userData.IsExist() {
		return nil, errs.Invalidated.Wrap(errs.ErrUserNotFound, "")
	}

	err = in.CheckPassword(userData.Password)
	if err != nil {
		return nil, err
	}

	token, err = userData.GenerateToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}
