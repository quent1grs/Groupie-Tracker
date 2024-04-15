package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func database() {
	user_id := 0
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	isEmpty, err := isUserTableEmpty(db)
	if err != nil {
		log.Fatal(err)
	}

	if isEmpty {
		user_id = 0
	} else {
		user_id, err = getNextID(db)
		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println("do you want to register ? (y/n)")
	var response string
	fmt.Scanln(&response)

	if response == "y" {
		register(db, user_id)
	}

	fmt.Println("do you want to delete all users ? (y/n)")
	fmt.Scanln(&response)

	if response == "y" {
		deleteAllUsers(db)
	}

	rows, err := db.Query("SELECT * FROM USER")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var pseudo string
		var email string
		var password string
		err = rows.Scan(&id, &pseudo, &email, &password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, pseudo, email, password)
	}
}

func register(db *sql.DB, id int) {
	if id > 1 {
		id = id - 1
	}
	pseudo := ""
	email := ""
	password := ""

	fmt.Println("Enter your pseudo: ")
	fmt.Scanln(&pseudo)
	fmt.Println("Enter your email: ")
	fmt.Scanln(&email)
	fmt.Println("Enter your password: ")
	fmt.Scanln(&password)

	id++

	stmt, err := db.Prepare("INSERT INTO USER(id, pseudo, email, password) VALUES(?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id, pseudo, email, password)

	if err != nil {
		log.Fatal(err)
	}
}

func getNextID(db *sql.DB) (int, error) {
	var maxID int
	err := db.QueryRow("SELECT MAX(id) FROM USER").Scan(&maxID)
	if err != nil {
		return 0, err
	}
	return maxID + 1, nil
}

func isUserTableEmpty(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM USER").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func deleteAllUsers(db *sql.DB) {
	stmt, err := db.Prepare("DELETE FROM USER")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
