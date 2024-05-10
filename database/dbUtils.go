package database

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// )

// func ResetStatus() {
// 	db, err := sql.Open("sqlite3", "./database/db.sqlite")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// 	_, err = db.Exec("UPDATE USER SET status = 'offline'")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("[EVENT] All users status reset to offline.")
// }
