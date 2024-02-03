package users

import (
	"github.com/golang-jwt/jwt"
	"test_ina_bank/config"
	"test_ina_bank/internal/common/errs"
	"test_ina_bank/pkg/utils"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

func (u *User) IsExist() bool {
	return u.Id != 0
}

func (u *User) HashPassword() error {
	hashPass, err := utils.HashPassword(u.Password)
	if err != nil {
		return errs.Failed.Wrap(errs.ErrGeneratePassword, "")
	}
	u.Password = hashPass
	return nil
}

func (u *User) GenerateToken() (*Token, error) {
	expire := time.Now().Add(config.Config.General.TokenDuration)
	claims := &JwtCustomClaims{
		UserId:   u.Id,
		Email:    u.Email,
		LoginAt:  time.Now(),
		ExpireAt: expire,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}

	// access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config.General.SecretKey))
	if err != nil {
		return nil, errs.Failed.New(err.Error())
	}

	// refresh token
	rtExpire := time.Now().Add(config.Config.General.RefreshTokenDuration)
	rtClaims := jwt.StandardClaims{
		ExpiresAt: rtExpire.Unix(),
		Subject:   u.Email,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rt, err := refreshToken.SignedString([]byte(config.Config.General.SecretKey))
	if err != nil {
		return nil, errs.Failed.New(err.Error())
	}

	return &Token{
		AccessToken:  t,
		RefreshToken: rt,
	}, nil
}

func (u *User) Replace(user *User) (err error) {
	*u = *user
	u.Password, err = utils.HashPassword(u.Password)
	if err != nil {
		return errs.ErrHashPassword
	}
	return nil
}

type UserId struct {
	Id int `json:"id" binding:"required"`
}

type ListUsers struct {
	Data []*User `json:"data"`
}
