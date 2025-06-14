package config

import "fmt"

// database queries
const (
	CheckDBQuery  = "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?"
	CreateDBQuery = "CREATE DATABASE IF NOT EXISTS %s"
	UseDBQuery    = "USE %s"
)

// mysql queries

var (
	DBPath = fmt.Sprintf("%s.%s", GetMysqlDBName(), GetMysqlTable())
)

const (
	CreateTableQuery   = "CREATE TABLE IF NOT EXISTS %s (id VARCHAR(36) UNIQUE NOT NULL PRIMARY KEY, name VARCHAR(36) NOT NULL, last_name VARCHAR(36) NOT NULL, username VARCHAR(36) UNIQUE NOT NULL, email VARCHAR(36) UNIQUE NOT NULL, password VARCHAR(100) NOT NULL)"
	CheckExistsQuery   = "SELECT 1 FROM %s WHERE username = ? LIMIT 1"
	GetByUsernameQuery = "SELECT id,name,last_name,username,email,password FROM %s WHERE username = ?"
	NewUserQuery       = "INSERT INTO %s (id,name,last_name,username,email,password) VALUES (?,?,?,?,?,?)"
	DeleteQuery        = "DELETE FROM %s WHERE username = ?"
	UpdateQuery        = "UPDATE %s SET name = ?, last_name = ?, email = ? WHERE username = ?"
	ChangePwdQuery     = "UPDATE %s SET password = ? WHERE username = ?"
)

// mysql test queries
const (
	CreateTableTest   = "CREATE TABLE IF NOT EXISTS table_name (id VARCHAR(36) UNIQUE NOT NULL PRIMARY KEY, name VARCHAR(36) NOT NULL, last_name VARCHAR(36) NOT NULL, username VARCHAR(36) UNIQUE NOT NULL, email VARCHAR(36) UNIQUE NOT NULL, password VARCHAR(100) NOT NULL)"
	CheckExistsTest   = "SELECT 1 FROM %s WHERE username = ? LIMIT 1"
	GetByUsernameTest = "SELECT id,name,last_name,username,email,password FROM  WHERE username = ?"
	NewUserTest       = "INSERT INTO  (id,name,last_name,username,email,password) VALUES (?,?,?,?,?,?)"
	DeleteUserTest    = "DELETE FROM WHERE username = ?"
	UpdateUserTest    = "UPDATE SET name = ?, last_name = ?, email = ? WHERE username = ?"
	ChangePwdTest     = "UPDATE SET password = ? WHERE username = ?"
)
