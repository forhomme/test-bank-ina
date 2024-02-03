package query

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/domain/users"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type ListUserHandler decorators.QueryHandler[*users.User, *users.ListUsers]

type listUserRepository struct {
	dbRepo users.QueryRepository
}

func NewListUserRepository(dbRepo users.QueryRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.QueryHandler[*users.User, *users.ListUsers] {
	return decorators.ApplyQueryDecorators[*users.User, *users.ListUsers](
		listUserRepository{dbRepo: dbRepo},
		log,
		tracer,
	)
}

func (l listUserRepository) Handle(ctx context.Context, _ *users.User) (*users.ListUsers, error) {
	usersData, err := l.dbRepo.ListUser(ctx)
	return &users.ListUsers{Data: usersData}, err
}
