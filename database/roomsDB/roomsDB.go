package roomsdb

import (
	"fmt"
	"log"

	db "groupietracker/database/DB_Connection"
)

func InsertRoomInDatabase(idRoom int, roomOwner int, maxPlayers int, name string, idGame int) {

	conn := db.GetDB()

	_, err := conn.Exec("INSERT INTO ROOMS (id, created_by, max_player, name, id_game) VALUES (?, ?, ?, ?, ?)", idRoom, roomOwner, maxPlayers, name, idGame)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] Room" + string(rune(idRoom)) + " created in database.")
}

func DeleteAllRooms() {
	conn := db.GetDB()

	_, err := conn.Exec("DELETE FROM ROOMS")
	if err != nil {
		log.Fatal(err)
	}
}

func GetNextAvailableID() int {

	conn := db.GetDB()

	rows, err := conn.Query("SELECT id FROM ROOMS")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var idRoom int
	for rows.Next() {
		err = rows.Scan(&idRoom)
		if err != nil {
			log.Fatal(err)
		}
	}
	return idRoom + 1
}
