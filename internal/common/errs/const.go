package errs

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrHashPassword     = errors.New("error hash password")
	ErrGeneratePassword = errors.New("cannot hashing password")
	ErrTaskNotFound     = errors.New("task not found")
	ErrorSigningMethod  = errors.New("unexpected signing method")
	ErrorTokenNotValid  = errors.New("token not valid")
	ErrAuthUserNotFound = errors.New("authorization user not found")
)
