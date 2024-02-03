package users

import (
	"github.com/golang-jwt/jwt"
	"test_ina_bank/config"
	"test_ina_bank/internal/common/errs"
	"time"
)

type RefreshToken struct {
	Email        string `json:"email" validate:"required,email"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (r *RefreshToken) ParsingToken(user *User) error {
	tokenRt, err := jwt.Parse(r.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.Invalidated.Wrap(errs.ErrorSigningMethod, "")
		}
		return []byte(config.Config.General.SecretKey), nil
	})
	if err != nil {
		return errs.Invalidated.Wrap(errs.ErrorTokenNotValid, "")
	}

	if rtClaims, ok := tokenRt.Claims.(jwt.MapClaims); ok && tokenRt.Valid {
		if rtClaims["sub"].(string) != user.Email {
			return errs.Invalidated.Wrap(errs.ErrorTokenNotValid, "")
		}
		if int64(rtClaims["exp"].(float64)) < time.Now().Unix() {
			return errs.Invalidated.Wrap(errs.ErrorTokenNotValid, "")
		}
		return nil
	}
	return errs.ErrorTokenNotValid
}
