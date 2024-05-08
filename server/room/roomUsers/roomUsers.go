package roomUsers

// Content of room_users de db.sqlite: id_room ; id_user ; score

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertUserInRoom(idRoom string, idUser string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO ROOMS_USERS (id_room, id_user) VALUES (?, ?)", idRoom, idUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] User" + idUser + " inserted in room.")
}

func DeleteUserFromRoom(idRoom string, idUser string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM ROOMS_USERS WHERE id_room = ? AND id_user = ?", idRoom, idUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] User" + idUser + " deleted from room.")
}

func AddScoreToUser(idRoom string, idUser string, score int) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE ROOMS_USERS SET score = score + ? WHERE id_room = ? AND id_user = ?", score, idRoom, idUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] Score added to user" + idUser + ".")
}

func GetScoreOfUser(idRoom string, idUser string) int {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var score int
	err = db.QueryRow("SELECT score FROM ROOMS_USERS WHERE id_room = ? AND id_user = ?", idRoom, idUser).Scan(&score)
	if err != nil {
		log.Fatal(err)
	}
	return score
}
