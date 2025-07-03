package auth

import (
	"fmt"
	"time"

	entity "go-manage-hex/internal/core/user"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	SecretKey string
	Duration  time.Duration
}

func NewJWTService(secret string, duration time.Duration) *JWTService {
	return &JWTService{
		SecretKey: secret,
		Duration:  duration,
	}
}

func (j *JWTService) GenerateJWT(username, password string) (string, error) {
	now := time.Now()

	claims := entity.Claims{
		Username: username,
		Password: password, // aunque usualmente no conviene poner password en el JWT!
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.Duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.SecretKey)) // usar j.SecretKey, no jwtKey global
}

func (j *JWTService) ValidateJWT(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &entity.Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inválido")
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("token inválido: %w", err)
	}

	claims, ok := token.Claims.(*entity.Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("token inválido")
	}

	return claims.Username, nil
}
