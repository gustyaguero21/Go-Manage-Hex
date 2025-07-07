package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (a *MockAuthService) GenerateJWT(username string) (string, error) {
	args := a.Called(username)
	return args.String(0), args.Error(1)
}

func (a *MockAuthService) ValidateJWT(token string) (string, error) {
	args := a.Called(token)
	return args.String(0), args.Error(1)
}

func TestRequireAuth(t *testing.T) {

	mockAuthService := new(MockAuthService)
	mw := NewMiddleware(mockAuthService)

	tests := []struct {
		Name           string
		AuthHeader     string
		MockValidate   func()
		ExpectedCode   int
		ExpectedBody   string
		ExpectUsername string
		ExpectNext     bool
	}{
		{
			Name:         "No Authorization header",
			AuthHeader:   "",
			MockValidate: func() {},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: `{"error":"authorization header missing"}`,
			ExpectNext:   false,
		},
		{
			Name:         "Malformed Authorization header",
			AuthHeader:   "Token abcdefg",
			MockValidate: func() {},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: `{"error":"authorization header format must be Bearer {token}"}`,
			ExpectNext:   false,
		},
		{
			Name:       "Invalid token",
			AuthHeader: "Bearer invalidtoken",
			MockValidate: func() {
				mockAuthService.On("ValidateJWT", "invalidtoken").Return("", errors.New("invalid token")).Once()
			},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: `{"error":"invalid token"}`,
			ExpectNext:   false,
		},
		{
			Name:       "Valid token",
			AuthHeader: "Bearer validtoken",
			MockValidate: func() {
				mockAuthService.On("ValidateJWT", "validtoken").Return("user123", nil).Once()
			},
			ExpectedCode:   http.StatusOK,
			ExpectNext:     true,
			ExpectUsername: "user123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mockAuthService.ExpectedCalls = nil
			tt.MockValidate()

			router := gin.New()
			router.Use(mw.RequireAuth)

			router.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"username": c.GetString("username"),
				})
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.AuthHeader != "" {
				req.Header.Set("Authorization", tt.AuthHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedCode, w.Code)
			if !tt.ExpectNext {
				assert.JSONEq(t, tt.ExpectedBody, w.Body.String())
			} else {
				assert.Contains(t, w.Body.String(), tt.ExpectUsername)
			}

			mockAuthService.AssertExpectations(t)
		})
	}
}
