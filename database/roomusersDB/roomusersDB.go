package roomusersdb

import (
	"fmt"
	db "groupietracker/database/DB_Connection"
	"log"
)

func InsertUserInRoomUsers(roomID int, userID int) {
	fmt.Println("[DEBUG] insertUserInRoomUsers() called.")
	fmt.Println("[DEBUG] roomusersdb.insertUserInRoomUsers() Room ID : " + string(rune(roomID)))

	conn := db.GetDB()

	_, err := conn.Exec("INSERT INTO ROOM_USERS (id_room, id_user) VALUES (?, ?)", roomID, userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] User" + string(rune(userID)) + " added to room" + string(rune(roomID)) + ".")
}

func ResetTable() {
	conn := db.GetDB()

	_, err := conn.Exec("DROP TABLE ROOM_USERS")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec("CREATE TABLE ROOM_USERS (id_room INTEGER, id_user INTEGER, score INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
}

func GetRoomAssociatedWithUser(username int) int {
	conn := db.GetDB()

	rows, err := conn.Query("SELECT id_room FROM ROOM_USERS WHERE id_user = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var roomID int
	for rows.Next() {
		err = rows.Scan(&roomID)
		if err != nil {
			log.Fatal(err)
		}
	}
	return roomID
}

func DeleteAllRoomUsers() {
	conn := db.GetDB()

	_, err := conn.Exec("DELETE FROM ROOM_USERS")
	if err != nil {
		log.Fatal(err)
	}
}
