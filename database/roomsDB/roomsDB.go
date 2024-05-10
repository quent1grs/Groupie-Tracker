package roomsdb

import (
	"fmt"
	"log"

	db "groupietracker/database/DB_Connection" // importez votre package DB_Connection
)

func InsertRoomInDatabase(idRoom int, roomOwner string, maxPlayers int, name string, idGame int) {
	fmt.Println("[DEBUG] insertRoomInDatabase() called.")
	defer fmt.Println("[DEBUG] insertRoomInDatabase() ended.")

	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("INSERT INTO ROOMS (id, created_by, max_player, name, id_game) VALUES (?, ?, ?, ?, ?)", idRoom, roomOwner, maxPlayers, name, idGame)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] Room" + string(idRoom) + " created in database.")
}

func DeleteAllRooms() {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("DELETE FROM ROOMS")
	if err != nil {
		log.Fatal(err)
	}
}

func GetNextAvailableID() int {
	fmt.Println("[DEBUG] roomsdb.GetNextAvailableID() called.")
	defer fmt.Println("[DEBUG] roomsdb.GetNextAvailableID() ended.")

	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

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
