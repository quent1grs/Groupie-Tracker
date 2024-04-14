package main

import (
	"groupietracker/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.Database()
}
