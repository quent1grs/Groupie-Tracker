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

func SetAsConnected(username string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE USER SET status = 'online' WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
}

func SetAsDisconnected(username string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE USER SET status = 'offline' WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
}

func DisconnectAllUsers() {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE USER SET status = 'offline'")
	if err != nil {
		log.Fatal(err)
	}
}
