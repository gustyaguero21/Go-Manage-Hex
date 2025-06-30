package user

import (
	"database/sql"
	"fmt"
	"go-manage-hex/cmd/config"
	mysqlrepo "go-manage-hex/internal/core/user"
)

type UserMysql struct {
	DB *sql.DB
}

func NewUserMysql(db *sql.DB) mysqlrepo.MysqlRepository {
	return &UserMysql{DB: db}
}

func (um *UserMysql) CreateTable(tableName string) error {
	query := fmt.Sprintf(config.CreateTableQuery, tableName)

	_, err := um.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserMysql) GetByUsername(username string) (mysqlrepo.User, error) {
	query := fmt.Sprintf(config.GetByUsernameQuery, config.GetMysqlTable())

	var user mysqlrepo.User

	err := um.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Name,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return mysqlrepo.User{}, config.ErrUserNotFound
	}

	return user, nil
}

func (um *UserMysql) NewUser(user mysqlrepo.User) error {
	query := fmt.Sprintf(config.NewUserQuery, config.GetMysqlTable())

	_, err := um.DB.Exec(query, user.ID, user.Name, user.LastName, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserMysql) DeleteUser(username string) error {
	query := fmt.Sprintf(config.DeleteQuery, config.GetMysqlTable())

	_, err := um.DB.Exec(query, username)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserMysql) UpdateUser(username string, user mysqlrepo.User) error {
	query := fmt.Sprintf(config.UpdateQuery, config.GetMysqlTable())

	_, err := um.DB.Exec(query, user.Name, user.LastName, user.Email, username)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserMysql) ChangePwd(newPwd, username string) error {
	query := fmt.Sprintf(config.ChangePwdQuery, config.GetMysqlTable())

	_, err := um.DB.Exec(query, newPwd, username)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserMysql) CheckExists(username string) bool {
	query := fmt.Sprintf(config.CheckExistsQuery, config.GetMysqlTable())

	var exists string

	return um.DB.QueryRow(query, username).Scan(&exists) == nil

}
