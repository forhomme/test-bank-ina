package app

import (
	"test_ina_bank/internal/usecase/command"
	"test_ina_bank/internal/usecase/query"
)

type TaskApplication struct {
	Commands TaskCommands
	Queries  TaskQueries
}

type TaskCommands struct {
	InsertTaskHandler command.InsertTaskHandler
	UpdateTaskHandler command.UpdateTaskHandler
	DeleteTaskHandler command.DeleteTaskHandler
}

type TaskQueries struct {
	ListTaskHandler query.ListTaskHandler
	GetTaskHandler  query.GetTaskHandler
}
