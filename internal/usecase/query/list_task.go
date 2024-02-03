package query

import (
	"context"
	"test_ina_bank/internal/common/decorators"
	"test_ina_bank/internal/domain/tasks"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

type ListTaskHandler decorators.QueryHandler[*tasks.Task, *tasks.ListTask]

type listTaskRepository struct {
	dbRepo tasks.QueryRepository
}

func NewListTaskRepository(dbRepo tasks.QueryRepository, log *baselogger.Logger, tracer *telemetry.OtelSdk) decorators.QueryHandler[*tasks.Task, *tasks.ListTask] {
	return decorators.ApplyQueryDecorators[*tasks.Task, *tasks.ListTask](
		listTaskRepository{dbRepo: dbRepo},
		log,
		tracer,
	)
}

func (l listTaskRepository) Handle(ctx context.Context, _ *tasks.Task) (*tasks.ListTask, error) {
	listData, err := l.dbRepo.ListTask(ctx)
	return &tasks.ListTask{Data: listData}, err
}
