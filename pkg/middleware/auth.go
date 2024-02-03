package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"test_ina_bank/config"
	"test_ina_bank/internal/common/errs"
	"test_ina_bank/pkg/utils"
)

type AuthUser struct {
	UserId int
	Email  string
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		if len(strings.Split(bearerToken, " ")) < 2 {
			c.JSON(http.StatusUnauthorized, utils.ParseMessage(errs.Unauthorized.New("missing authorization token")))
			c.Abort()
			return
		}
		tokenString := strings.Split(bearerToken, " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				err := errs.Unauthorized.New(errs.ErrorSigningMethod.Error())
				return nil, err
			}
			return []byte(config.Config.General.SecretKey), nil
		})
		if err != nil {
			c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
		if !ok {
			err = errs.Unauthorized.Wrap(errs.ErrorTokenNotValid, "")
			c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
			return
		}

		err = claims.Valid()
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ParseMessage(errs.Unauthorized.Wrap(err, "")))
			c.Abort()
			return
		}

		authUser := &AuthUser{
			UserId: int(claims["user_id"].(float64)),
			Email:  claims["email"].(string),
		}
		c.Set("authUser", authUser)
		c.Next()
	}
}

func GetAuthUser(c *gin.Context) (*AuthUser, error) {
	value, ok := c.Get("authUser")
	if !ok {
		return nil, errs.Unauthorized.New(errs.ErrAuthUserNotFound.Error())
	}
	authUser, ok := value.(*AuthUser)
	if !ok {
		return nil, errs.Unauthorized.New(errs.ErrAuthUserNotFound.Error())
	}
	return authUser, nil
}
