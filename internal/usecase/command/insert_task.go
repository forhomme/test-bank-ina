package command

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/domain/tasks"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type InsertTaskHandler decorators.CommandHandler[*tasks.Task]

type insertTaskRepository struct {
	dbRepo tasks.CommandRepository
}

func NewInsertTaskRepository(dbRepo tasks.CommandRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.CommandHandler[*tasks.Task] {
	return decorators.ApplyCommandDecorators[*tasks.Task](
		insertTaskRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (i insertTaskRepository) Handle(ctx context.Context, in *tasks.Task) (err error) {
	err = i.dbRepo.InsertTask(ctx, in)
	return err
}
