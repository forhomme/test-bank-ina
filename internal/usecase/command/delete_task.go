package command

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/domain/tasks"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type DeleteTaskHandler decorators.CommandHandler[*tasks.TaskId]

type deleteTaskRepository struct {
	dbRepo tasks.CommandRepository
}

func NewDeleteTaskRepository(dbRepo tasks.CommandRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.CommandHandler[*tasks.TaskId] {
	return decorators.ApplyCommandDecorators[*tasks.TaskId](
		deleteTaskRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (d deleteTaskRepository) Handle(ctx context.Context, in *tasks.TaskId) (err error) {
	err = d.dbRepo.DeleteTask(ctx, in.Id)
	return err
}
