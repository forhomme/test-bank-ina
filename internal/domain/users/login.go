package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"test_ina_bank/internal/common/errs"
)

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

var WrongUsernameOrPasswordError = errors.New("wrong username or password")

func (l *Login) CheckPassword(hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(l.Password))
	if err != nil {
		return errs.Invalidated.Wrap(WrongUsernameOrPasswordError, "")
	}
	return nil
}
