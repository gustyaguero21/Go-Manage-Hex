package middleware

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}
