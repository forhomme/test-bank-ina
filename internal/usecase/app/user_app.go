package app

import (
	"test_ina_bank/internal/usecase/command"
	"test_ina_bank/internal/usecase/query"
)

type UserApplication struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	InsertUserHandler command.InsertUserHandler
	UpdateUserHandler command.UpdateUserHandler
	DeleteUserHandler command.DeleteUserHandler
}

type Queries struct {
	ListUserHandler     query.ListUserHandler
	GetUserHandler      query.GetUserHandler
	LoginHandler        query.LoginHandler
	RefreshTokenHandler query.RefreshTokenHandler
}
