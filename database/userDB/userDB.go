package userdb

import (
	db "groupietracker/database/DB_Connection"
	"log"
)

func GetIDFromUsername(username string) int {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	var id int
	err := conn.QueryRow("SELECT id FROM USER WHERE username = ?", username).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func SetAsConnected(username string) {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("UPDATE USER SET status = 'online' WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
}

func SetAsDisconnected(username string) {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("UPDATE USER SET status = 'offline' WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
}

func DisconnectAllUsers() {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection

	_, err := conn.Exec("UPDATE USER SET status = 'offline'")
	if err != nil {
		log.Fatal(err)
	}
}
