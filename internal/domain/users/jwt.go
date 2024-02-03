package users

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtCustomClaims struct {
	UserId   int       `json:"user_id"`
	Email    string    `json:"email"`
	LoginAt  time.Time `json:"login_at"`
	ExpireAt time.Time `json:"expire_at"`
	jwt.StandardClaims
}
