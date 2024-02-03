package command

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/domain/users"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type DeleteUserHandler decorators.CommandHandler[*users.UserId]

type deleteUserRepository struct {
	dbRepo users.CommandRepository
}

func NewDeleteUserRepository(dbRepo users.CommandRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.CommandHandler[*users.UserId] {
	return decorators.ApplyCommandDecorators[*users.UserId](
		deleteUserRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (d deleteUserRepository) Handle(ctx context.Context, in *users.UserId) (err error) {
	err = d.dbRepo.DeleteUser(ctx, in.Id)
	return err
}
