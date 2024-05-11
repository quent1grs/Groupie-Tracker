package userdb

import (
	"fmt"
	db "groupietracker/database/DB_Connection"
	"log"
)

func GetIDFromUsername(username string) int {
	conn := db.GetDB()
	fmt.Println("[DEBUG] userdb.GetIDFromUsername() called.")
	fmt.Println("[DEBUG] userdb.GetIDFromUsername() username : " + username)
	fmt.Println("[DEBUG] userdb.GetIDFromUsername() conn : ")
	fmt.Println(conn)

	rows, err := conn.Query("SELECT id FROM USER WHERE username = ?", username)
	if err != nil {
		fmt.Println("[DEBUG] userdb.GetIDFromUsername() : err = ")
		fmt.Println(err)
		log.Fatal(err)
	}

	defer rows.Close()

	var id int
	for rows.Next() {
		fmt.Println("[DEBUG] userdb.GetIDFromUsername() : rows.Next() called.")
		fmt.Println("[DEBUG] userdb.GetIDFromUsername() : rows = ")
		fmt.Println(rows)
		err = rows.Scan(&id)
		fmt.Println("[DEBUG] userdb.GetIDFromUsername() : id = " + string(rune(id)))
		if err != nil {
			fmt.Println("[DEBUG] userdb.GetIDFromUsername() : err = ")
			fmt.Println(err)
			log.Fatal(err)
		}
		break // Add a break statement to exit the loop after the first iteration
	}

	fmt.Println("[DEBUG] id : " + string(rune(id)))
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
