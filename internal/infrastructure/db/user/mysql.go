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

func (mr *UserMysql) CreateTable(tableName string) error {
	query := fmt.Sprintf(config.CreateTableQuery, tableName)

	_, err := mr.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table. Error: %s", err)
	}
	return nil
}
