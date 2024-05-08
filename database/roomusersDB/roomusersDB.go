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
