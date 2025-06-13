package user

import (
	"context"
	"fmt"
	entity "go-manage-hex/internal/core/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMysqlRepository struct {
	CreateTableFn   func(tableName string) error
	GetByUsernameFn func(username string) (entity.User, error)
	CheckExistsFn   func(username string) bool
	NewUserFn       func(user entity.User) error
	DeleteUserFn    func(username string) error
	UpdateUserFn    func(username string, user entity.User) error
	ChangePwdFn     func(newPwd, username string) error
}

func (m *mockMysqlRepository) CreateTable(tableName string) error {
	if m.CreateTableFn != nil {
		return m.CreateTableFn(tableName)
	}
	return nil
}

func (m *mockMysqlRepository) GetByUsername(username string) (entity.User, error) {
	if m.GetByUsernameFn != nil {
		return m.GetByUsernameFn(username)
	}
	return entity.User{}, nil
}

func (m *mockMysqlRepository) CheckExists(username string) bool {
	if m.CheckExistsFn != nil {
		return m.CheckExistsFn(username)
	}
	return false
}

func (m *mockMysqlRepository) NewUser(user entity.User) error {
	if m.NewUserFn != nil {
		return m.NewUserFn(user)
	}
	return nil
}

func (m *mockMysqlRepository) DeleteUser(username string) error {
	if m.DeleteUserFn != nil {
		return m.DeleteUserFn(username)
	}
	return nil
}

func (m *mockMysqlRepository) UpdateUser(username string, user entity.User) error {
	if m.UpdateUserFn != nil {
		return m.UpdateUserFn(username, user)
	}
	return nil
}

func (m *mockMysqlRepository) ChangePwd(newPwd, username string) error {
	if m.ChangePwdFn != nil {
		return m.ChangePwdFn(newPwd, username)
	}
	return nil
}

func TestSearchUser(t *testing.T) {
	test := []struct {
		Name         string
		Username     string
		mockExists   bool
		mockUser     entity.User
		mockGetErr   error
		expectedUser entity.User
	}{
		{
			Name:       "SearchUser_Success",
			Username:   "johndoe",
			mockExists: true,
			mockUser: entity.User{
				ID:       "1",
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234.",
			},
		},
		{
			Name:       "SearchUser_NotFound",
			Username:   "johncito",
			mockExists: false,
			mockUser:   entity.User{},
			mockGetErr: fmt.Errorf("user not found"),
		},
		{
			Name:       "SearchUser_Error",
			Username:   "johncito",
			mockExists: true,
			mockUser:   entity.User{},
			mockGetErr: fmt.Errorf("error searching user. Error: "),
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			mockRepo := mockMysqlRepository{
				CheckExistsFn: func(username string) bool {
					return tt.mockExists
				},
				GetByUsernameFn: func(username string) (entity.User, error) {
					return tt.mockUser, tt.mockGetErr
				},
			}
			service := NewUserService(&mockRepo)

			found, err := service.SearchUser(context.Background(), tt.Username)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if found.Username != tt.Username {
				assert.NotEqual(t, found.Username, tt.Username)
			} else {
				assert.Equal(t, found.Username, tt.Username)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	test := []struct {
		Name        string
		mockExists  bool
		mockGetErr  error
		User        entity.User
		expectedErr error
	}{
		{
			Name:       "CreateUser_Success",
			mockExists: false,
			User: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234567",
			},
			expectedErr: nil,
		},
		{
			Name:        "CreateUser_ErrUserAlreadyExists",
			mockExists:  true,
			expectedErr: fmt.Errorf("user already exists"),
		},
		{
			Name:        "CreateUser_ErrInvalidEmail",
			mockExists:  false,
			expectedErr: fmt.Errorf("invalid email address"),
		},
		{
			Name:       "CreateUser_ErrInvalidPassword",
			mockExists: false,
			User: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234567.",
			},
			expectedErr: fmt.Errorf("invalid password"),
		},
		{
			Name:       "CreateUser_Err",
			mockExists: false,
			User: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234567",
			},
			expectedErr: fmt.Errorf("error creating user. Error: some db error"),
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			mockRepo := mockMysqlRepository{
				CheckExistsFn: func(username string) bool {
					return tt.mockExists
				},
				NewUserFn: func(user entity.User) error {
					return tt.expectedErr
				},
			}
			service := NewUserService(&mockRepo)

			_, err := service.CreateUser(context.Background(), tt.User)

			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
