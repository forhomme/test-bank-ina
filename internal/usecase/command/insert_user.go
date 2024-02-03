package command

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/domain/users"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type InsertUserHandler decorators.CommandHandler[*users.User]

type insertUserRepository struct {
	dbRepo users.CommandRepository
}

func NewInsertUserRepository(dbRepo users.CommandRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.CommandHandler[*users.User] {
	return decorators.ApplyCommandDecorators[*users.User](
		insertUserRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (i insertUserRepository) Handle(ctx context.Context, in *users.User) (err error) {
	err = in.HashPassword()
	if err != nil {
		return err
	}
	err = i.dbRepo.InsertUser(ctx, in)
	return err
}
