package command

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/domain/users"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type UpdateUserHandler decorators.CommandHandler[*users.User]

type updateUserRepository struct {
	dbRepo users.CommandRepository
}

func NewUpdateUserRepository(dbRepo users.CommandRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.CommandHandler[*users.User] {
	return decorators.ApplyCommandDecorators[*users.User](
		updateUserRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (u updateUserRepository) Handle(ctx context.Context, in *users.User) (err error) {
	return u.dbRepo.UpdateUser(
		ctx,
		in.Id,
		func(u *users.User) (*users.User, error) {
			err = u.Replace(in)
			if err != nil {
				return nil, err
			}
			return u, nil
		})
}
