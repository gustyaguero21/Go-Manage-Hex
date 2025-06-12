package user

import (
	"context"
	"fmt"

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
		return mysqlUser.User{}, fmt.Errorf("user not found")
	}

	search, searchErr := us.Repo.GetByUsername(username)
	if searchErr != nil {
		return mysqlUser.User{}, searchErr
	}

	return search, nil
}

func (us *UserServices) CreateUser(ctx context.Context, user mysqlUser.User) (created mysqlUser.User, err error) {
	if us.Repo.CheckExists(user.Username) {
		return mysqlUser.User{}, fmt.Errorf("user already exists")
	}

	uID := uuid.NewString()

	if !validator.ValidateEmail(user.Email) {
		return mysqlUser.User{}, fmt.Errorf("invalid email address")
	}

	if !validator.ValidatePassword(user.Password) {
		return mysqlUser.User{}, fmt.Errorf("invalid password")
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
		return fmt.Errorf("user not found")
	}

	if deleteErr := us.Repo.DeleteUser(username); deleteErr != nil {
		return fmt.Errorf("error deleting user. Error: %s", deleteErr)
	}

	return nil
}

func (us *UserServices) UpdateUser(ctx context.Context, username string, user mysqlUser.User) (updated mysqlUser.User, err error) {
	if !us.Repo.CheckExists(username) {
		return mysqlUser.User{}, fmt.Errorf("user not found")
	}

	if !validator.ValidateEmail(user.Email) {
		return mysqlUser.User{}, fmt.Errorf("invalid email address")
	}

	if updateErr := us.Repo.UpdateUser(username, user); updateErr != nil {
		return mysqlUser.User{}, fmt.Errorf("error updating user. Error: %s", updateErr)
	}

	return user, nil
}

func (us *UserServices) ChangeUserPwd(ctx context.Context, newPwd, username string) error {
	if !us.Repo.CheckExists(username) {
		return fmt.Errorf("user not found")
	}

	if !validator.ValidatePassword(newPwd) {
		return fmt.Errorf("invalid password")
	}

	if changePwd := us.Repo.ChangePwd(newPwd, username); changePwd != nil {
		return fmt.Errorf("error changing password. Error: %s", changePwd)
	}

	return nil
}
