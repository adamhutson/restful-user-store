package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error

	host := "localhost"
	port := 5432
	user := "adam"
	password := "password"
	database := "restful-user-store"

	dbInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	if err := createTables(); err != nil {
		panic(err)
	}
}

func createTables() error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (FirstName TEXT, LastName TEXT, UserID TEXT PRIMARY KEY NOT NULL)")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS groups (GroupName TEXT PRIMARY KEY NOT NULL)")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS userGroups (UserID TEXT NOT NULL, GroupName TEXT NOT NULL, PRIMARY KEY(UserID, GroupName))")
	if err != nil {
		return err
	}
	return nil
}
