package middleware

import (
	"net/http"
	"strings"

	auth "go-manage-hex/internal/core/user"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	AuthService auth.Authorization
}

func NewMiddleware(authSvc auth.Authorization) *Middleware {
	return &Middleware{AuthService: authSvc}
}

func (m *Middleware) RequireAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
		c.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer {token}"})
		c.Abort()
		return
	}

	tokenString := parts[1]

	username, err := m.AuthService.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
		c.Abort()
		return
	}

	c.Set("username", username)
	c.Next()
}
