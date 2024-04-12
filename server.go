package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gameID := 1

	gameName := "Game 1"

	stmt, err := db.Prepare("INSERT INTO GAMES(id, name) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(gameID, gameName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Values inserted successfully")
}
