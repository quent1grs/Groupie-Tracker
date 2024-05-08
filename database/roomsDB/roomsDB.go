package roomsdb

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertRoomInDatabase(idRoom int, roomOwner string, maxPlayers int, name string, idGame int) {
	fmt.Println("[DEBUG] insertRoomInDatabase() called.")
	defer fmt.Println("[DEBUG] insertRoomInDatabase() ended.")

	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO ROOMS (id, created_by, max_player, name, id_game) VALUES (?, ?, ?, ?, ?)", idRoom, roomOwner, maxPlayers, name, idGame)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] Room" + string(idRoom) + " created in database.")
}

func DeleteAllRooms() {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM ROOMS")
	if err != nil {
		log.Fatal(err)
	}
}

func GetNextAvailableID() int {
	fmt.Println("[DEBUG] roomsdb.GetNextAvailableID() called.")
	defer fmt.Println("[DEBUG] roomsdb.GetNextAvailableID() ended.")

	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM ROOMS")
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
