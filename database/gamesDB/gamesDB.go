package gamesdb

import (
	"fmt"
	"log"

	db "groupietracker/database/DB_Connection" // importez votre package DB_Connection
)

func CreateGame(name string) int {
	fmt.Println("[DEBUG] CreateGame() called.")
	defer fmt.Println("[DEBUG] CreateGame() ended.")

	fmt.Println("[DEBUG] name : " + name)

	id := nextAvailableID()

	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("INSERT INTO GAMES (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[EVENT] Game" + string(rune(id)) + " created.")
	return id
}

func DeleteGameFromDB(gameID string) {
	fmt.Println("[DEBUG] DeleteGameFromDB() called.")

	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("DELETE FROM GAMES WHERE id = ?", gameID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] Game" + gameID + " deleted from database.")
	fmt.Println("[DEBUG] DeleteGameFromDB() ended.")
}

func ResetTable() {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("DELETE FROM GAMES")
	if err != nil {
		log.Fatal(err)
	}
}

func EnumerateGames() {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

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

func nextAvailableID() int {
	fmt.Println("[DEBUG] nextAvailableID() called.")
	defer fmt.Println("[DEBUG] nextAvailableID() ended.")

	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	rows, err := conn.Query("SELECT id FROM GAMES")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	i := 1
	for rows.Next() {
		var id int
		rows.Scan(&id)
		if id != i {
			return i
		}
		i++
	}
	return 0
}
