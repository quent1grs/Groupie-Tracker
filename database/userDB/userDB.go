package userdb

import (
	"fmt"
	db "groupietracker/database/DB_Connection"
	"log"
)

func GetIDFromUsername(username string) int {
	conn := db.GetDB()
	fmt.Println(conn)

	rows, err := conn.Query("SELECT id FROM USER WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	defer rows.Close()

	var id int
	for rows.Next() {
		fmt.Println(rows)
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		break
	}

	return id
}

func SetAsConnected(username string) {
	conn := db.GetDB()

	_, err := conn.Exec("UPDATE USER SET status = 'online' WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
}

func SetAsDisconnected(username string) {
	conn := db.GetDB()

	_, err := conn.Exec("UPDATE USER SET status = 'offline' WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
}

func DisconnectAllUsers() {
	conn := db.GetDB()

	_, err := conn.Exec("UPDATE USER SET status = 'offline'")
	if err != nil {
		log.Fatal(err)
	}
}
