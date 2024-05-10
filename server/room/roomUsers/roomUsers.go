package roomUsers

import (
	"fmt"
	db "groupietracker/database/DB_Connection"
	"log"
)

func InsertUserInRoom(idRoom string, idUser string) {
	_, err := db.Conn.Exec("INSERT INTO ROOMS_USERS (id_room, id_user) VALUES (?, ?)", idRoom, idUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] User" + idUser + " inserted in room.")
}

func DeleteUserFromRoom(idRoom string, idUser string) {
	_, err := db.Conn.Exec("DELETE FROM ROOMS_USERS WHERE id_room = ? AND id_user = ?", idRoom, idUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] User" + idUser + " deleted from room.")
}

func AddScoreToUser(idRoom string, idUser string, score int) {
	_, err := db.Conn.Exec("UPDATE ROOMS_USERS SET score = score + ? WHERE id_room = ? AND id_user = ?", score, idRoom, idUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] Score added to user" + idUser + ".")
}

func GetScoreOfUser(idRoom string, idUser string) int {
	var score int
	err := db.Conn.QueryRow("SELECT score FROM ROOMS_USERS WHERE id_room = ? AND id_user = ?", idRoom, idUser).Scan(&score)
	if err != nil {
		log.Fatal(err)
	}
	return score
}
