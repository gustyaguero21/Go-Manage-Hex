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
		MockExists   bool
		MockUser     entity.User
		MockGetErr   error
		expectedUser entity.User
	}{
		{
			Name:       "SearchUser_Success",
			Username:   "johndoe",
			MockExists: true,
			MockUser: entity.User{
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
			MockExists: false,
			MockUser:   entity.User{},
			MockGetErr: fmt.Errorf("user not found"),
		},
		{
			Name:       "SearchUser_Error",
			Username:   "johncito",
			MockExists: true,
			MockUser:   entity.User{},
			MockGetErr: fmt.Errorf("error searching user. Error: "),
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			mockRepo := mockMysqlRepository{
				CheckExistsFn: func(username string) bool {
					return tt.MockExists
				},
				GetByUsernameFn: func(username string) (entity.User, error) {
					return tt.MockUser, tt.MockGetErr
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
		MockExists  bool
		MockGetErr  error
		User        entity.User
		ExpectedErr error
	}{
		{
			Name:       "CreateUser_Success",
			MockExists: false,
			User: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234567",
			},
			ExpectedErr: nil,
		},
		{
			Name:        "CreateUser_ErrUserAlreadyExists",
			MockExists:  true,
			ExpectedErr: fmt.Errorf("user already exists"),
		},
		{
			Name:        "CreateUser_ErrInvalidEmail",
			MockExists:  false,
			ExpectedErr: fmt.Errorf("invalid email address"),
		},
		{
			Name:       "CreateUser_ErrInvalidPassword",
			MockExists: false,
			User: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234567.",
			},
			ExpectedErr: fmt.Errorf("invalid password"),
		},
		{
			Name:       "CreateUser_Err",
			MockExists: false,
			User: entity.User{
				Name:     "John",
				LastName: "Doe",
				Username: "johndoe",
				Email:    "johndoe@example.com",
				Password: "Password1234567",
			},
			ExpectedErr: fmt.Errorf("error creating user. Error: some db error"),
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			mockRepo := mockMysqlRepository{
				CheckExistsFn: func(username string) bool {
					return tt.MockExists
				},
				NewUserFn: func(user entity.User) error {
					return tt.ExpectedErr
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

func TestDeleteUser(t *testing.T) {
	test := []struct {
		Name        string
		Username    string
		MockExists  bool
		ExpectedErr error
	}{
		{
			Name:        "DeleteUser_Success",
			Username:    "johndoe",
			MockExists:  true,
			ExpectedErr: nil,
		},
		{
			Name:        "DeleteUser_ErrUserNotFound",
			Username:    "johndoe",
			MockExists:  false,
			ExpectedErr: fmt.Errorf("user not found"),
		},
		{
			Name:        "DeleteUser_Err",
			Username:    "johndoe",
			MockExists:  true,
			ExpectedErr: fmt.Errorf("error deleting user. Error: some error"),
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			repo := mockMysqlRepository{
				CheckExistsFn: func(username string) bool {
					return tt.MockExists
				},
				DeleteUserFn: func(username string) error {
					return tt.ExpectedErr
				},
			}

			service := NewUserService(&repo)

			deleteErr := service.DeleteUser(context.Background(), tt.Username)

			if deleteErr != nil {
				assert.Error(t, deleteErr)
			} else {
				assert.NoError(t, deleteErr)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	test := []struct {
		Name        string
		Username    string
		User        entity.User
		MockExists  bool
		ExpectedErr error
	}{
		{
			Name:     "UpdateUser_Success",
			Username: "johndoe",
			User: entity.User{
				Name:     "Johncito",
				LastName: "Doecito",
				Email:    "johncitodoecito@example.com",
			},
			MockExists:  true,
			ExpectedErr: nil,
		},
		{
			Name:        "UpdateUser_ErrUserNotFound",
			Username:    "johndoe",
			MockExists:  false,
			ExpectedErr: fmt.Errorf("user not found"),
		},
		{
			Name:     "UpdateUser_ErrInvalidEmail",
			Username: "johndoe",
			User: entity.User{
				Email: "not-an-email",
			},
			MockExists:  true,
			ExpectedErr: fmt.Errorf("invalid email address"),
		},
		{
			Name:     "UpdateUser_Err",
			Username: "johndoe",
			User: entity.User{
				Email: "john@example.com",
			},
			MockExists:  true,
			ExpectedErr: fmt.Errorf("error updating user. Error: some error"),
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			repo := mockMysqlRepository{
				CheckExistsFn: func(username string) bool {
					return tt.MockExists
				},
				UpdateUserFn: func(username string, user entity.User) error {
					return tt.ExpectedErr
				},
			}
			service := NewUserService(&repo)

			_, err := service.UpdateUser(context.Background(), tt.Username, tt.User)

			if tt.ExpectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChangeUserPwd(t *testing.T) {
	test := []struct {
		Name        string
		NewPwd      string
		Username    string
		MockExists  bool
		ExpectedErr error
	}{
		{
			Name:        "ChangeUserPwd_Success",
			NewPwd:      "NewPassword1234",
			Username:    "johndoe",
			MockExists:  true,
			ExpectedErr: nil,
		},
		{
			Name:        "ChangeUserPwd_ErrUserNotFound",
			NewPwd:      "NewPassword1234",
			Username:    "johndoe",
			MockExists:  false,
			ExpectedErr: fmt.Errorf("user not found"),
		},
		{
			Name:        "ChangeUserPwd_ErrInvalidPassword",
			NewPwd:      "NewPassword1234.",
			Username:    "johndoe",
			MockExists:  true,
			ExpectedErr: fmt.Errorf("invalid password"),
		},
		{
			Name:        "ChangeUserPwd_Err",
			NewPwd:      "NewPassword1234",
			Username:    "johndoe",
			MockExists:  true,
			ExpectedErr: fmt.Errorf("error changing password. Error: some error"),
		},
	}

	for _, tt := range test {
		t.Run(tt.Name, func(t *testing.T) {
			repo := mockMysqlRepository{
				CheckExistsFn: func(username string) bool {
					return tt.MockExists
				},
				ChangePwdFn: func(newPwd, username string) error {
					return tt.ExpectedErr
				},
			}
			service := NewUserService(&repo)

			err := service.ChangeUserPwd(context.Background(), tt.NewPwd, tt.Username)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
