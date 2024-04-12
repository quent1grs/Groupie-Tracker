package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:5500)/index")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func createTable() {
	db := connect()

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS USER(id INTEGER PRIMARY KEY, pseudo TEXT NOT NULL, email TEXT NOT NULL, password TEXT NOT NULL);")

	checkErr(err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS ROOMS (id INTEGER PRIMARY KEY, created_by INTEGER NOT NULL, max_player INTEGER NOT NULL, name TEXT NOT NULL, id_game INTEGER, FOREIGN KEY (created_by) REFERENCES USER(id), FOREIGN KEY (id_game) REFERENCES GAMES(id));")

	checkErr(err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS ROOM_USERS (id_room INTEGER, id_user INTEGER, score INTEGER, FOREIGN KEY (id_room) REFERENCES ROOMS(id), FOREIGN KEY (id_user) REFERENCES USER(id), PRIMARY KEY (id_room, id_user));")

	checkErr(err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS GAMES (id INTEGER PRIMARY KEY, name TEXT NOT NULL);")

	checkErr(err)

}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
