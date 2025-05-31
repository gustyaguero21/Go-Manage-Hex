package config

// database queries
const (
	CheckDBQuery  = "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?"
	CreateDBQuery = "CREATE DATABASE IF NOT EXISTS %s"
	UseDBQuery    = "USE %s"
)

// mysql queries
const (
	CreateTableQuery = "CREATE TABLE IF NOT EXISTS %s (id VARCHAR(36) UNIQUE NOT NULL PRIMARY KEY, name VARCHAR(36) NOT NULL, last_name VARCHAR(36) NOT NULL, username VARCHAR(36) UNIQUE NOT NULL, email VARCHAR(36) UNIQUE NOT NULL, password VARCHAR(36) NOT NULL)"
)
