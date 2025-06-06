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
	CreateTableQuery = "CREATE TABLE IF NOT EXISTS %s (id VARCHAR(36) UNIQUE NOT NULL PRIMARY KEY, name VARCHAR(36) NOT NULL, last_name VARCHAR(36) NOT NULL, username VARCHAR(36) UNIQUE NOT NULL, email VARCHAR(36) UNIQUE NOT NULL, password VARCHAR(36) NOT NULL)"
	GetByNameQuery   = "SELECT id,name,last_name,username,email,password FROM %s WHERE name = ?"
	NewUserQuery     = "INSERT INTO %s (id,name,last_name,username,email,password) VALUES (?,?,?,?,?,?)"
	DeleteQuery      = "DELETE FROM %s WHERE name = ?"
	UpdateQuery      = "UPDATE %s SET name = ?, last_name = ?, email = ? WHERE name = ?"
	ChangePwdQuery   = "UPDATE %s SET password = ? WHERE name = ?"
)

// mysql test queries
const (
	CreateTableTest = "CREATE TABLE IF NOT EXISTS table_name (id VARCHAR(36) UNIQUE NOT NULL PRIMARY KEY, name VARCHAR(36) NOT NULL, last_name VARCHAR(36) NOT NULL, username VARCHAR(36) UNIQUE NOT NULL, email VARCHAR(36) UNIQUE NOT NULL, password VARCHAR(36) NOT NULL)"
	GetByNameTest   = "SELECT id,name,last_name,username,email,password FROM  WHERE name = ?"
	NewUserTest     = "INSERT INTO  (id,name,last_name,username,email,password) VALUES (?,?,?,?,?,?)"
	DeleteUserTest  = "DELETE FROM WHERE name = ?"
	UpdateUserTest  = "UPDATE SET name = ?, last_name = ?, email = ? WHERE name = ?"
	ChangePwdTest   = "UPDATE SET password = ? WHERE name = ?"
)
