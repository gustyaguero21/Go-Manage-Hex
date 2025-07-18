package user

import (
	"context"
	"errors"
	"fmt"
	entity "go-manage-hex/internal/core/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUsecases struct {
	mock.Mock
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockUsecases) SearchUser(ctx context.Context, username string) (entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUsecases) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUsecases) DeleteUser(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}

func (m *MockUsecases) UpdateUser(ctx context.Context, username string, user entity.User) (entity.User, error) {
	args := m.Called(ctx, username, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUsecases) ChangeUserPwd(ctx context.Context, newPwd, username string) error {
	args := m.Called(ctx, newPwd, username)
	return args.Error(0)
}

func (m *MockUsecases) Login(ctx context.Context, username, password string) error {
	args := m.Called(ctx, username, password)
	return args.Error(0)
}

func (a *MockAuthService) GenerateJWT(username string) (string, error) {
	args := a.Called(username)
	return args.String(0), args.Error(1)
}

func (a *MockAuthService) ValidateJWT(token string) (string, error) {
	args := a.Called(token)
	return args.String(0), args.Error(1)
}

func TestSearchUserHandler(t *testing.T) {
	mockUsecase := new(MockUsecases)
	handler := UserHandler{Service: mockUsecase}

	tests := []struct {
		Name           string
		Username       string
		ExpectedUser   entity.User
		MockFunc       func()
		ExpectedStatus int
	}{
		{
			Name:     "Success",
			Username: "johndoe",
			ExpectedUser: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			MockFunc: func() {
				mockUsecase.
					On("SearchUser", mock.Anything, "johndoe").
					Return(entity.User{}, nil).
					Once()
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:     "Invalid Query Param",
			Username: "",
			ExpectedUser: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			MockFunc: func() {
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:     "Error",
			Username: "johndoe",
			ExpectedUser: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			MockFunc: func() {
				mockUsecase.
					On("SearchUser", mock.Anything, "johndoe").
					Return(entity.User{}, errors.New("error searching user")).
					Once()
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockFunc()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			reqURL := fmt.Sprintf("/search?username=%s", tt.Username)
			c.Request = httptest.NewRequest(http.MethodGet, reqURL, nil)
			c.Request.Header.Set("Content-Type", "application/json")

			handler.SearchUserHandler(c)

			assert.Equal(t, tt.ExpectedStatus, w.Code)
		})
	}
}

func TestCreateUserHandler(t *testing.T) {
	mockUsecase := new(MockUsecases)
	handler := &UserHandler{Service: mockUsecase}

	tests := []struct {
		Name           string
		Body           string
		ExpectedUser   entity.User
		MockFunc       func()
		ExpectedStatus int
	}{
		{
			Name: "Success",
			Body: `{"name":"John","last_name":"Doe","username":"johndoe","email":"johndoe@example.com","password":"Password1234"}`,
			ExpectedUser: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			MockFunc: func() {
				mockUsecase.
					On("CreateUser", mock.Anything, mock.AnythingOfType("user.User")).
					Return(entity.User{}, nil).
					Once()
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "Invalid JSON",
			Body:           `{"name":"John", "password": "Password1234}`,
			MockFunc:       func() {},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "Error",
			Body: `{"name":"John","last_name":"Doe","username":"johndoe","email":"johndoe@example.com","password":"Password1234"}`,
			ExpectedUser: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234",
			},
			MockFunc: func() {
				mockUsecase.
					On("CreateUser", mock.Anything, mock.AnythingOfType("user.User")).
					Return(entity.User{}, errors.New("error creating user")).
					Once()
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.MockFunc()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tt.Body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateUserHandler(c)

			assert.Equal(t, tt.ExpectedStatus, w.Code)
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	mockUsecase := new(MockUsecases)
	handler := UserHandler{Service: mockUsecase}

	tests := []struct {
		Name           string
		Username       string
		Confirmation   string
		MockFunc       func()
		ExpectedStatus int
	}{
		{
			Name:         "Success",
			Username:     "johndoe",
			Confirmation: "True",
			MockFunc: func() {
				mockUsecase.
					On("DeleteUser", mock.Anything, "johndoe").
					Return(nil).
					Once()
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:         "Invalid Query Params",
			Username:     "",
			Confirmation: "True",
			MockFunc: func() {
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:         "Confirmation Error",
			Username:     "johndoe",
			Confirmation: "",
			MockFunc: func() {
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:         "Confirmation False",
			Username:     "johndoe",
			Confirmation: "False",
			MockFunc: func() {
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:         "Error",
			Username:     "johndoe",
			Confirmation: "True",
			MockFunc: func() {
				mockUsecase.
					On("DeleteUser", mock.Anything, "johndoe").
					Return(fmt.Errorf("error deleting user")).
					Once()
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.MockFunc()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			reqURL := fmt.Sprintf("/delete?username=%s&confirmation=%s", tt.Username, tt.Confirmation)
			c.Request = httptest.NewRequest(http.MethodDelete, reqURL, nil)
			c.Request.Header.Set("Content-Type", "application/json")

			handler.DeleteUserHandler(c)

			assert.Equal(t, tt.ExpectedStatus, w.Code)
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	mockUsecase := new(MockUsecases)
	handler := UserHandler{Service: mockUsecase}

	tests := []struct {
		Name           string
		Username       string
		Update         string
		MockFunc       func()
		ExpectedStatus int
	}{
		{
			Name:     "Success",
			Username: "johndoe",
			Update:   `{"name":"Johncito","last_name":"Doecito","email":"johncitodoecito@example.com"}`,
			MockFunc: func() {
				mockUsecase.
					On("UpdateUser", mock.Anything, "johndoe", entity.User{
						Name:     "Johncito",
						LastName: "Doecito",
						Email:    "johncitodoecito@example.com",
					}).
					Return(entity.User{}, nil).
					Once()
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:     "Invalid Query Param",
			Username: "",
			Update:   `{"name":"Johncito","last_name":"Doecito","email":"johncitodoecito@example.com"}`,
			MockFunc: func() {
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:     "Invalid JSON",
			Username: "johndoe",
			Update:   `{"name":"Johncito","last_name":"Doecito","email":"johncitodoecito@example.com}`,
			MockFunc: func() {
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:     "Invalid Body Content",
			Username: "johndoe",
			Update:   `{"name":"","last_name":"Doecito","email":"johncitodoecito@example.com"}`,
			MockFunc: func() {
			},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:     "Error",
			Username: "johndoe",
			Update:   `{"name":"Johncito","last_name":"Doecito","email":"johncitodoecito@example.com"}`,
			MockFunc: func() {
				mockUsecase.
					On("UpdateUser", mock.Anything, "johndoe", entity.User{
						Name:     "Johncito",
						LastName: "Doecito",
						Email:    "johncitodoecito@example.com",
					}).
					Return(entity.User{}, fmt.Errorf("error updating user")).
					Once()
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.MockFunc()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyReader := strings.NewReader(tt.Update)
			reqURL := fmt.Sprintf("/update?username=%s", tt.Username)

			c.Request = httptest.NewRequest(http.MethodPatch, reqURL, bodyReader)
			c.Request.Header.Set("Content-Type", "application/json")

			handler.UpdateUserHandler(c)

			assert.Equal(t, tt.ExpectedStatus, w.Code)
		})
	}
}

func TestChangePwdHandler(t *testing.T) {
	mockUsecase := new(MockUsecases)
	handler := UserHandler{Service: mockUsecase}

	tests := []struct {
		Name           string
		Body           string
		MockFunc       func()
		ExpectedStatus int
	}{
		{
			Name: "Success",
			Body: `{"username":"johndoe","new_pwd":"newpassword"}`,
			MockFunc: func() {
				mockUsecase.
					On("ChangeUserPwd", mock.Anything, "newpassword", "johndoe").
					Return(nil).
					Once()
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "Invalid JSON Body",
			Body:           `{"username":"",`,
			MockFunc:       func() {},
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name: "Error",
			Body: `{"username":"johndoe","new_pwd":"newpassword"}`,
			MockFunc: func() {
				mockUsecase.
					On("ChangeUserPwd", mock.Anything, mock.Anything, mock.Anything).
					Return(fmt.Errorf("error changing password")).
					Once()
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.MockFunc()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPatch, "/change-password", strings.NewReader(tt.Body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.ChangePwdHandler(c)

			assert.Equal(t, tt.ExpectedStatus, w.Code)
		})
	}
}

func TestLoginUser(t *testing.T) {
	mockUsecase := new(MockUsecases)
	mockAuthService := new(MockAuthService)
	handler := UserHandler{
		Service:     mockUsecase,
		AuthService: mockAuthService,
	}

	tests := []struct {
		Name            string
		Login           string
		MockLogin       func()
		MockGenerateJWT func()
		ExpectedStatus  int
	}{
		{
			Name:  "success",
			Login: `{"username": "john", "password": "doe123"}`,
			MockLogin: func() {
				mockUsecase.On("Login", mock.Anything, "john", "doe123").Return(nil).Once()
			},
			MockGenerateJWT: func() {
				mockAuthService.On("GenerateJWT", "john").Return("mocked-token", nil).Once()
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:            "invalid json",
			Login:           `{"username": "john"`,
			MockLogin:       func() {},
			MockGenerateJWT: func() {},
			ExpectedStatus:  http.StatusBadRequest,
		},
		{
			Name:  "invalid credentials",
			Login: `{"username": "john", "password": "wrong"}`,
			MockLogin: func() {
				mockUsecase.On("Login", mock.Anything, "john", "wrong").Return(errors.New("unauthorized")).Once()
			},
			MockGenerateJWT: func() {},
			ExpectedStatus:  http.StatusUnauthorized,
		},
		{
			Name:  "error generating jwt",
			Login: `{"username": "john", "password": "doe123"}`,
			MockLogin: func() {
				mockUsecase.On("Login", mock.Anything, "john", "doe123").Return(nil).Once()
			},
			MockGenerateJWT: func() {
				mockAuthService.On("GenerateJWT", "john").Return("", errors.New("internal")).Once()
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			tt.MockLogin()
			tt.MockGenerateJWT()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyReader := strings.NewReader(tt.Login)

			c.Request = httptest.NewRequest(http.MethodPost, "/login", bodyReader)
			c.Request.Header.Set("Content-Type", "application/json")

			handler.LoginUser(c)

			assert.Equal(t, tt.ExpectedStatus, w.Code)

		})
	}
}
