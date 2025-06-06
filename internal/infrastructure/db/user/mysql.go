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

func (um *UserMysql) GetByName(name string) (mysqlrepo.User, error) {
	query := fmt.Sprintf(config.GetByNameQuery, config.GetMysqlTable())

	var user mysqlrepo.User

	err := um.DB.QueryRow(query, name).Scan(
		&user.ID,
		&user.Name,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return mysqlrepo.User{}, sql.ErrNoRows
	}
	if err != nil {
		return mysqlrepo.User{}, err
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

func (um *UserMysql) DeleteUser(name string) error {
	query := fmt.Sprintf(config.DeleteQuery, config.GetMysqlTable())

	_, err := um.DB.Exec(query, name)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserMysql) UpdateUser(name string, user mysqlrepo.User) error {
	query := fmt.Sprintf(config.UpdateQuery, config.GetMysqlTable())

	_, err := um.DB.Exec(query, user.Name, user.LastName, user.Email, name)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserMysql) ChangePwd(newPwd, name string) error {
	query := fmt.Sprintf(config.ChangePwdQuery, config.GetMysqlTable())

	_, err := um.DB.Exec(query, newPwd, name)
	if err != nil {
		return err
	}
	return nil
}
