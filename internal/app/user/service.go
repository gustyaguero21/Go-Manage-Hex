package user

import (
	"context"

	mysqlUser "go-manage-hex/internal/core/user"
)

type UserServices struct {
	Repo mysqlUser.MysqlRepository
}

func NewUserService(repo mysqlUser.MysqlRepository) Usecases {
	return &UserServices{Repo: repo}
}

func (us *UserServices) CreateUser(ctx context.Context, user mysqlUser.User) (created mysqlUser.User, err error) {
	return mysqlUser.User{}, nil
}
