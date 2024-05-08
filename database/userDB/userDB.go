package userdb

import (
	"database/sql"
	"log"
)

func GetIDFromUsername(username string) int {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var id int
	err = db.QueryRow("SELECT id FROM USER WHERE username = ?", username).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}
