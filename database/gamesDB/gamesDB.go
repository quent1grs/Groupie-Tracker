package gamesdb

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateGame(name string) int {
	fmt.Println("[DEBUG] CreateGame() called.")
	defer fmt.Println("[DEBUG] CreateGame() ended.")

	fmt.Println("[DEBUG] name : " + name)

	id := nextAvailableID()

	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO GAMES (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[EVENT] Game" + string(rune(id)) + " created.")
	return id
}

func DeleteGameFromDB(gameID string) {
	fmt.Println("[DEBUG] DeleteGameFromDB() called.")

	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM GAMES WHERE id = ?", gameID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] Game" + gameID + " deleted from database.")
	fmt.Println("[DEBUG] DeleteGameFromDB() ended.")
}

func ResetTable() {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM GAMES")
	if err != nil {
		log.Fatal(err)
	}
}

func EnumerateGames() {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM GAMES")
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

	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM GAMES")
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
