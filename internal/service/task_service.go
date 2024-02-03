package service

import (
	"test_ina_bank/internal/adapters"
	"test_ina_bank/internal/usecase/app"
	"test_ina_bank/internal/usecase/command"
	"test_ina_bank/internal/usecase/query"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/sqlhandler"
	"test_ina_bank/pkg/telemetry"
)

func NewTaskService(log *baselogger.Logger, sqlHandler sqlhandler.SqlHandler, tracer *telemetry.OtelSdk) app.TaskApplication {
	mysqlRepo := adapters.NewMysqlRepository(log, sqlHandler, tracer)
	return app.TaskApplication{
		Commands: app.TaskCommands{
			InsertTaskHandler: command.NewInsertTaskRepository(mysqlRepo, log, tracer),
			UpdateTaskHandler: command.NewUpdateTaskRepository(mysqlRepo, log, tracer),
			DeleteTaskHandler: command.NewDeleteTaskRepository(mysqlRepo, log, tracer),
		},
		Queries: app.TaskQueries{
			ListTaskHandler: query.NewListTaskRepository(mysqlRepo, log, tracer),
			GetTaskHandler:  query.NewGetTaskRepository(mysqlRepo, log, tracer),
		},
	}
}
