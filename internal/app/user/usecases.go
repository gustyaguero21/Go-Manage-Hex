package user

import (
	"context"
	mysqlUser "go-manage-hex/internal/core/user"
)

type Usecases interface {
	CreateUser(ctx context.Context, user mysqlUser.User) (created mysqlUser.User, err error)
}
