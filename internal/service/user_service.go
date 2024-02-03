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

func NewUserService(log *baselogger.Logger, sqlHandler sqlhandler.SqlHandler, tracer *telemetry.OtelSdk) app.UserApplication {
	mysqlRepo := adapters.NewMysqlRepository(log, sqlHandler, tracer)
	return app.UserApplication{
		Commands: app.Commands{
			InsertUserHandler: command.NewInsertUserRepository(mysqlRepo, log, tracer),
			UpdateUserHandler: command.NewUpdateUserRepository(mysqlRepo, log, tracer),
			DeleteUserHandler: command.NewDeleteUserRepository(mysqlRepo, log, tracer),
		},
		Queries: app.Queries{
			ListUserHandler:     query.NewListUserRepository(mysqlRepo, log, tracer),
			GetUserHandler:      query.NewGetUserRepository(mysqlRepo, log, tracer),
			LoginHandler:        query.NewLoginRepository(mysqlRepo, log, tracer),
			RefreshTokenHandler: query.NewRefreshTokenRepository(mysqlRepo, log, tracer),
		},
	}
}
