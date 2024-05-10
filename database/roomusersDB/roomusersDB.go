package roomusersdb

import (
	"fmt"
	db "groupietracker/database/DB_Connection"
	"log"
)

func InsertUserInRoomUsers(roomID int, userID int) {
	fmt.Println("[DEBUG] insertUserInRoomUsers() called.")

	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("INSERT INTO ROOM_USERS (id_room, id_user) VALUES (?, ?)", roomID, userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] User" + string(userID) + " added to room" + string(roomID) + ".")
}

func ResetTable() {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("DROP TABLE ROOM_USERS")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec("CREATE TABLE ROOM_USERS (id_room INTEGER, id_user INTEGER, score INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteAllRoomUsers() {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("DELETE FROM ROOM_USERS")
	if err != nil {
		log.Fatal(err)
	}
}
