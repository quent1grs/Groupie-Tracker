package gamesdb

import (
	"fmt"
	"log"
	"math/rand"

	db "groupietracker/database/DB_Connection"
)

var digits = []rune("0123456789")

func CreateGame(name string, gametype string) int {

	id := getRandomId()

	conn := db.GetDB()

	_, err := conn.Exec("INSERT INTO GAMES (id, name, gameMode) VALUES (?, ?, ?)", id, name, gametype)
	if err != nil {
		log.Fatal(err)
	}

	return id
}

func DeleteGameFromDB(gameID string) {

	conn := db.GetDB()

	_, err := conn.Exec("DELETE FROM GAMES WHERE id = ?", gameID)
	if err != nil {
		log.Fatal(err)
	}
}

func ResetTable() {
	conn := db.GetDB()

	_, err := conn.Exec("DELETE FROM GAMES")
	if err != nil {
		log.Fatal(err)
	}
}

func EnumerateGames() {
	conn := db.GetDB()

	rows, err := conn.Query("SELECT id, name FROM GAMES")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println("Game", id, ":", name)
	}
}

func getRandomId() int {
	id := int(digits[rand.Intn(len(digits))])
	for !isIdUnique(id) {
		id = 0
		for i := 0; i < 6; i++ {
			id = id*10 + int(digits[rand.Intn(len(digits))])
		}
	}
	return id
}

func GetGameMode(gameID int) string {
	conn := db.GetDB()

	rows, err := conn.Query("SELECT gameMode FROM GAMES WHERE id = ?", gameID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var mode string
	for rows.Next() {
		err = rows.Scan(&mode)
		if err != nil {
			log.Fatal(err)
		}
	}
	return mode
}

func isIdUnique(id int) bool {
	conn := db.GetDB()

	rows, err := conn.Query("SELECT id FROM GAMES WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		return false
	}
	return true
}
