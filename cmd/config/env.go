package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if loadErr := godotenv.Load(); loadErr != nil {
		log.Fatal(loadErr)
	}
	log.Print(".env file loaded successfully")
}

func GetMysqlUser() string {
	return os.Getenv("MYSQL_USER")
}

func GetMysqlPwd() string {
	return os.Getenv("MYSQL_ROOT_PASSWORD")
}

func GetMysqlDBHost() string {
	return os.Getenv("MYSQL_DB_HOST")
}

func GetMysqlDBPort() string {
	return os.Getenv("MYSQL_DB_PORT")
}

func GetMysqlDBName() string {
	return os.Getenv("MYSQL_DB_NAME")
}

func GetMysqlTable() string {
	return os.Getenv("MYSQL_TABLE_NAME")
}

func GetDSNRoot() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", GetMysqlUser(), GetMysqlPwd(), GetMysqlDBHost(), GetMysqlDBPort())

	return dsn
}

func GetDSN_DB() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", GetMysqlUser(), GetMysqlPwd(), GetMysqlDBHost(), GetMysqlDBPort(), GetMysqlDBName())

	return dsn
}
