package tasks

import "context"

type CommandRepository interface {
	InsertTask(ctx context.Context, in *Task) (err error)
	UpdateTask(ctx context.Context, id int, updateFn func(u *Task) (*Task, error)) error
	DeleteTask(ctx context.Context, id int) error
}

type QueryRepository interface {
	ListTask(ctx context.Context) ([]*TaskModel, error)
	GetTaskById(ctx context.Context, id int) (*TaskModel, error)
}
