package query

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/common/errs"
	"test_ina_bank/internal/domain/tasks"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type GetTaskHandler decorators.QueryHandler[*tasks.TaskId, *tasks.TaskModel]

type getTaskRepository struct {
	dbRepo tasks.QueryRepository
}

func NewGetTaskRepository(dbRepo tasks.QueryRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.QueryHandler[*tasks.TaskId, *tasks.TaskModel] {
	return decorators.ApplyQueryDecorators[*tasks.TaskId, *tasks.TaskModel](
		getTaskRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (g getTaskRepository) Handle(ctx context.Context, in *tasks.TaskId) (*tasks.TaskModel, error) {
	taskData, err := g.dbRepo.GetTaskById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if !taskData.IsExist() {
		return nil, errs.ErrTaskNotFound
	}
	return taskData, nil
}
