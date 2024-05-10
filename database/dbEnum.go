package database

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// )

// func EnumerateConnectedUsers() {
// 	db, err := sql.Open("sqlite3", "./database/db.sqlite")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// 	rows, err := db.Query("SELECT * FROM USER WHERE status = 'connected'")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id int
// 		var username string
// 		var email string
// 		var password string
// 		var status string
// 		var sessioncookie string
// 		err = rows.Scan(&id, &username, &email, &password, &status, &sessioncookie)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("ID: %d, Username: %s, Cookie: %s\n", id, username, sessioncookie)
// 	}
// }

// func EnumerateUsersTable() {
// 	db, err := sql.Open("sqlite3", "./database/db.sqlite")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// 	rows, err := db.Query("SELECT * FROM USER")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id int
// 		var username string
// 		var email string
// 		var password string
// 		var status string
// 		var sessioncookie string
// 		err = rows.Scan(&id, &username, &email, &password, &status, &sessioncookie)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("ID: %d, Username: %s, Email: %s, Password: %s, Status: %s, Cookie: %s\n", id, username, email, password, status, sessioncookie)
// 	}
// }

// func EnumerateRoomsTable() {
// 	db, err := sql.Open("sqlite3", "./database/db.sqlite")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// 	rows, err := db.Query("SELECT * FROM ROOMS")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id string
// 		var roomOwner string
// 		err = rows.Scan(&id, &roomOwner)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("ID: %s, RoomOwner: %s\n", id, roomOwner)
// 	}
// }

// func EnumerateRoomsUsersTable() {
// 	db, err := sql.Open("sqlite3", "./database/db.sqlite")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// 	rows, err := db.Query("SELECT * FROM ROOMS_USERS")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var roomID string
// 		var username string
// 		err = rows.Scan(&roomID, &username)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("RoomID: %s, Username: %s\n", roomID, username)
// 	}
// }

// func EnumerateUsersGamesTable() {
// 	db, err := sql.Open("sqlite3", "./database/db.sqlite")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// 	rows, err := db.Query("SELECT * FROM USERS_GAMES")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var username string
// 		var gameID string
// 		err = rows.Scan(&username, &gameID)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("Username: %s, GameID: %s\n", username, gameID)
// 	}
// }
