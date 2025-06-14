package user

import (
	"context"
	"fmt"

	"go-manage-hex/cmd/config"
	mysqlUser "go-manage-hex/internal/core/user"

	"github.com/google/uuid"
	"github.com/gustyaguero21/go-core/pkg/encrypter"
	"github.com/gustyaguero21/go-core/pkg/validator"
)

type UserServices struct {
	Repo mysqlUser.MysqlRepository
}

func NewUserService(repo mysqlUser.MysqlRepository) Usecases {
	return &UserServices{Repo: repo}
}

func (us *UserServices) SearchUser(ctx context.Context, username string) (search mysqlUser.User, err error) {
	if !us.Repo.CheckExists(username) {
		return mysqlUser.User{}, config.ErrUserNotFound
	}

	search, searchErr := us.Repo.GetByUsername(username)
	if searchErr != nil {
		return mysqlUser.User{}, fmt.Errorf("error searching user. Error: %s", searchErr)
	}

	return search, nil
}

func (us *UserServices) CreateUser(ctx context.Context, user mysqlUser.User) (created mysqlUser.User, err error) {
	if us.Repo.CheckExists(user.Username) {
		return mysqlUser.User{}, config.ErrUserAlreadyExists
	}

	uID := uuid.NewString()

	if !validator.ValidateEmail(user.Email) {
		return mysqlUser.User{}, config.ErrInvalidEmail
	}

	if !validator.ValidatePassword(user.Password) {
		return mysqlUser.User{}, config.ErrInvalidPassword
	}

	hash, hashErr := encrypter.PasswordEncrypter(user.Password)
	if hashErr != nil {
		return mysqlUser.User{}, hashErr
	}

	user.ID = uID
	user.Password = string(hash)

	if createErr := us.Repo.NewUser(user); createErr != nil {
		return mysqlUser.User{}, fmt.Errorf("error creating user. Error: %s", createErr)
	}

	return user, nil
}

func (us *UserServices) DeleteUser(ctx context.Context, username string) error {
	if !us.Repo.CheckExists(username) {
		return config.ErrUserNotFound
	}

	if deleteErr := us.Repo.DeleteUser(username); deleteErr != nil {
		return fmt.Errorf("error deleting user. Error: %s", deleteErr)
	}

	return nil
}

func (us *UserServices) UpdateUser(ctx context.Context, username string, user mysqlUser.User) (updated mysqlUser.User, err error) {
	if !us.Repo.CheckExists(username) {
		return mysqlUser.User{}, config.ErrUserNotFound
	}

	if !validator.ValidateEmail(user.Email) {
		return mysqlUser.User{}, config.ErrInvalidEmail
	}

	if updateErr := us.Repo.UpdateUser(username, user); updateErr != nil {
		return mysqlUser.User{}, fmt.Errorf("error updating user. Error: %s", updateErr)
	}

	return user, nil
}

func (us *UserServices) ChangeUserPwd(ctx context.Context, newPwd, username string) error {
	if !us.Repo.CheckExists(username) {
		return config.ErrUserNotFound
	}

	if !validator.ValidatePassword(newPwd) {
		return config.ErrInvalidPassword
	}

	if changePwd := us.Repo.ChangePwd(newPwd, username); changePwd != nil {
		return fmt.Errorf("error changing password. Error: %s", changePwd)
	}

	return nil
}
