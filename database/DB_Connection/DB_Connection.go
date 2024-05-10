package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB

func init() {
	var err error
	Conn, err = sql.Open("sqlite3", "database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *sql.DB {
	return Conn
}
