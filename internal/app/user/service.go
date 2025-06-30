package user

import (
	"context"

	"go-manage-hex/cmd/config"
	mysqlUser "go-manage-hex/internal/core/user"

	"github.com/google/uuid"
	"github.com/gustyaguero21/go-core/pkg/apperror"
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
		return mysqlUser.User{}, apperror.AppError(config.ErrSearchingUser, config.ErrUserNotFound)
	}

	search, searchErr := us.Repo.GetByUsername(username)
	if searchErr != nil {
		return mysqlUser.User{}, apperror.AppError(config.ErrSearchingUser, searchErr)
	}

	return search, nil
}

func (us *UserServices) CreateUser(ctx context.Context, user mysqlUser.User) (created mysqlUser.User, err error) {
	if us.Repo.CheckExists(user.Username) {
		return mysqlUser.User{}, apperror.AppError(config.ErrCreatingUser, config.ErrUserAlreadyExists)
	}

	uID := uuid.NewString()

	if !validator.ValidateEmail(user.Email) {
		return mysqlUser.User{}, apperror.AppError(config.ErrCreatingUser, config.ErrInvalidEmail)
	}

	if !validator.ValidatePassword(user.Password) {
		return mysqlUser.User{}, apperror.AppError(config.ErrCreatingUser, config.ErrInvalidPassword)
	}

	hash, _ := encrypter.PasswordEncrypter(user.Password)

	user.ID = uID
	user.Password = string(hash)

	if createErr := us.Repo.NewUser(user); createErr != nil {
		return mysqlUser.User{}, apperror.AppError(config.ErrCreatingUser, createErr)
	}

	return user, nil
}

func (us *UserServices) DeleteUser(ctx context.Context, username string) error {
	if !us.Repo.CheckExists(username) {
		return apperror.AppError(config.ErrDeletingUser, config.ErrUserNotFound)
	}

	if deleteErr := us.Repo.DeleteUser(username); deleteErr != nil {
		return apperror.AppError(config.ErrDeletingUser, deleteErr)
	}

	return nil
}

func (us *UserServices) UpdateUser(ctx context.Context, username string, user mysqlUser.User) (updated mysqlUser.User, err error) {
	if !us.Repo.CheckExists(username) {
		return mysqlUser.User{}, apperror.AppError(config.ErrUpdatingUser, config.ErrUserNotFound)
	}

	if !validator.ValidateEmail(user.Email) {
		return mysqlUser.User{}, apperror.AppError(config.ErrUpdatingUser, config.ErrInvalidEmail)
	}

	if updateErr := us.Repo.UpdateUser(username, user); updateErr != nil {
		return mysqlUser.User{}, apperror.AppError(config.ErrUpdatingUser, updateErr)
	}

	return user, nil
}

func (us *UserServices) ChangeUserPwd(ctx context.Context, newPwd, username string) error {
	if !us.Repo.CheckExists(username) {
		return apperror.AppError(config.ErrChangingPwd, config.ErrUserNotFound)
	}

	if !validator.ValidatePassword(newPwd) {
		return apperror.AppError(config.ErrChangingPwd, config.ErrInvalidPassword)
	}

	hash, _ := encrypter.PasswordEncrypter(newPwd)

	if changePwdErr := us.Repo.ChangePwd(string(hash), username); changePwdErr != nil {
		return apperror.AppError(config.ErrChangingPwd, changePwdErr)
	}

	return nil
}
