package db

import (
	"database/sql"
	"fmt"
	"go-manage-hex/cmd/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func DatabaseConn() (*sql.DB, error) {
	tempDB, err := sql.Open("mysql", config.GetDSNRoot())
	if err != nil {
		return nil, err
	}

	if pingErr := tempDB.Ping(); pingErr != nil {
		log.Fatal("error pinging to database")
	}

	if err := checkExists(tempDB, config.GetMysqlDBName()); err != nil {
		log.Print("DATABASE NOT FOUND. CREATING....")
		if createErr := createDB(tempDB, config.GetMysqlDBName()); createErr != nil {
			return nil, fmt.Errorf("failed to create database: %w", createErr)
		}
		log.Print("DATABASE CREATED SUCCESSFULLY")
	} else {
		log.Print("DATABASE FOUND. USING")
	}

	db, err := sql.Open("mysql", config.GetDSN_DB())
	if err != nil {
		log.Fatal(err)
	}
	if useErr := useDB(db, config.GetMysqlDBName()); useErr != nil {
		log.Fatal(err)
	}

	return db, nil
}

func checkExists(db *sql.DB, dbName string) error {
	var exists string
	err := db.QueryRow(config.CheckDBQuery, dbName).Scan(&exists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("database not found")
	}
	return nil
}

func createDB(db *sql.DB, dbName string) error {
	query := fmt.Sprintf(config.CreateDBQuery, dbName)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func useDB(db *sql.DB, dbName string) error {
	query := fmt.Sprintf(config.UseDBQuery, dbName)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
