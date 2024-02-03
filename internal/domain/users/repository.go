package users

import "context"

type CommandRepository interface {
	InsertUser(ctx context.Context, in *User) (err error)
	UpdateUser(ctx context.Context, id int, updateFn func(u *User) (*User, error)) error
	DeleteUser(ctx context.Context, id int) error
}

type QueryRepository interface {
	ListUser(ctx context.Context) ([]*User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
