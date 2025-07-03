package user

import (
	"context"
	mysqlUser "go-manage-hex/internal/core/user"
)

type Usecases interface {
	SearchUser(ctx context.Context, username string) (search mysqlUser.User, err error)
	CreateUser(ctx context.Context, user mysqlUser.User) (created mysqlUser.User, err error)
	DeleteUser(ctx context.Context, username string) error
	UpdateUser(ctx context.Context, username string, user mysqlUser.User) (updated mysqlUser.User, err error)
	ChangeUserPwd(ctx context.Context, newPwd, username string) error
	Login(ctx context.Context, username, password string) error
}
