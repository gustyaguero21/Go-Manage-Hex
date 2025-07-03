package auth

import (
	"fmt"
	"time"

	claim "go-manage-hex/internal/infrastructure/http/middleware"

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

	claims := claim.Claims{
		Username: username,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.Duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTService) ValidateJWT(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claim.Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token sign")
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*claim.Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Username, nil
}
