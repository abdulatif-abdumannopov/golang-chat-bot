package models

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "./project.db")
	if err != nil {
		log.Fatal(err)
	}

	userTable := `
	CREATE TABLE IF NOT EXISTS users (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
		telegram integer
	)
	`
	_, err = db.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
