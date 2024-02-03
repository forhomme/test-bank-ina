package query

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/common/errs"
	"test_ina_bank/internal/domain/users"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type GetUserHandler decorators.QueryHandler[*users.UserId, *users.User]

type getUserRepository struct {
	dbRepo users.QueryRepository
}

func NewGetUserRepository(dbRepo users.QueryRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.QueryHandler[*users.UserId, *users.User] {
	return decorators.ApplyQueryDecorators[*users.UserId, *users.User](
		getUserRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (g getUserRepository) Handle(ctx context.Context, in *users.UserId) (*users.User, error) {
	userData, err := g.dbRepo.GetUserById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if !userData.IsExist() {
		return nil, errs.ErrUserNotFound
	}
	return userData, nil
}
