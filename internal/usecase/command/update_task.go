package command

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/domain/tasks"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type UpdateTaskHandler decorators.CommandHandler[*tasks.Task]

type updateTaskRepository struct {
	dbRepo tasks.CommandRepository
}

func NewUpdateTaskRepository(dbRepo tasks.CommandRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.CommandHandler[*tasks.Task] {
	return decorators.ApplyCommandDecorators[*tasks.Task](
		updateTaskRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (u updateTaskRepository) Handle(ctx context.Context, in *tasks.Task) (err error) {
	return u.dbRepo.UpdateTask(
		ctx,
		in.Id,
		func(t *tasks.Task) (*tasks.Task, error) {
			t.Replace(in)
			return t, nil
		})
}
