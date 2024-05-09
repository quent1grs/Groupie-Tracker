package roomusersdb

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertUserInRoomUsers(roomID int, userID int) {
	fmt.Println("[DEBUG] insertUserInRoomUsers() called.")

	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO ROOM_USERS (id_room, id_user) VALUES (?, ?)", roomID, userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] User" + string(userID) + " added to room" + string(roomID) + ".")
}

func ResetTable() {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("DROP TABLE ROOM_USERS")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE ROOM_USERS (id_room INTEGER, id_user INTEGER, score INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
}
