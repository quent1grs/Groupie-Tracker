package database

import (
	"database/sql"
	"fmt"
	"log"
)

func ShowUserDetails(username string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM USER WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		var email string
		var password string
		var status string
		var sessioncookie string
		err = rows.Scan(&id, &username, &email, &password, &status, &sessioncookie)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Username: %s, Email: %s, Password: %s, Status: %s, Cookie: %s\n", id, username, email, password, status, sessioncookie)
	}
}

func ShowRoomDetails(roomname string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM ROOMS WHERE id = ?", roomname)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var roomOwner string
		err = rows.Scan(&id, &roomOwner)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s, RoomOwner: %s\n", id, roomOwner)
	}
}
