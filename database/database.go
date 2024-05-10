package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	db "groupietracker/database/DB_Connection"
	"log"
)

func Database() {
	db := db.GetDB()
	GetUsers(db)
	isEmpty, err := isUserTableEmpty(db)
	if err != nil {
		log.Fatal(err)
	}

	if isEmpty {
		// user_id = 0
	} else {
		// user_id, err = getNextID(db)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func InsertFormData(username, email, password string) error {
	db := db.GetDB()

	stmt, err := db.Prepare("INSERT INTO USER(id, username, email, password, status, sessioncookie) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("erreur lors de la préparation de la requête d'insertion: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(nextAvailableID(), email, username, password, "offline", "")
	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion dans la base de données: %v", err)
	}

	fmt.Println("Données du formulaire insérées avec succès dans la base de données.")

	return nil
}

func Hash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func GetUsers(db *sql.DB) {
	rows, err := db.Query("SELECT id, username, email, password FROM USER")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, email, password string
		err = rows.Scan(&id, &username, &email, &password)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, username, email, password)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func nextAvailableID() int {
	db := db.GetDB()

	var maxID int
	err := db.QueryRow("SELECT MAX(id) FROM USER").Scan(&maxID)
	if err != nil {
		return 0
	}

	return maxID + 1
}

func isUserTableEmpty(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM USER").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func IsUsernameInDB(username string) bool {
	db := db.GetDB()

	rows, err := db.Query("SELECT username FROM USER WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return rows.Next()
}

func IsEmailInDB(email string) bool {
	db := db.GetDB()

	rows, err := db.Query("SELECT email FROM USER WHERE email = ?", email)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return rows.Next()
}

func IsPasswordCorrect(identifier string, password string) bool {
	if identifier == "" || password == "" {
		return false
	}
	if !IsUsernameInDB(identifier) && !IsEmailInDB(identifier) {
		return false
	}

	db := db.GetDB()

	var hashedPassword string
	var idOfRow int
	err := db.QueryRow("SELECT password FROM USER WHERE username = ? OR email = ?", identifier, identifier).Scan(&hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow("SELECT id FROM USER WHERE username = ? OR email = ?", identifier, identifier).Scan(&idOfRow)
	if err != nil {
		log.Fatal(err)
	}

	if !isIdentifierPresentInTheRow(identifier, idOfRow) {
		return false
	}

	var comparedPassword string
	err = db.QueryRow("SELECT password FROM USER WHERE id = ?", idOfRow).Scan(&comparedPassword)
	if err != nil {
		log.Fatal(err)
	}

	return Hash(password) == comparedPassword
}

func isIdentifierPresentInTheRow(identifier string, id int) bool {
	db := db.GetDB()

	var comparedUsername string
	var comparedEmail string
	err := db.QueryRow("SELECT username FROM USER WHERE id = ?", id).Scan(&comparedUsername)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow("SELECT email FROM USER WHERE id = ?", id).Scan(&comparedEmail)
	if err != nil {
		log.Fatal(err)
	}
	return comparedUsername == identifier || comparedEmail == identifier
}

func ResetSessionData() {
	db := db.GetDB()
	_, err := db.Exec("UPDATE USER SET status = 'offline'")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("UPDATE USER SET sessioncookie = ''")
	if err != nil {
		log.Fatal(err)
	}
}

// Nettoie la base de données des informations temporaires (ROOM_USERS, GAMES, ROOMS)
func InitDatabase() {
	db := db.GetDB()
	_, _ = db.Exec("DELETE FROM ROOM_USERS")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	_, _ = db.Exec("DELETE FROM GAMES")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	_, _ = db.Exec("DELETE FROM ROOMS")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	_, _ = db.Exec("UPDATE USER SET status = 'offline'")
	_, _ = db.Exec("UPDATE USER SET sessioncookie = ''")
	fmt.Println("[EVENT:database.InitDatabase()] ROOM_USERS, GAMES, ROOMS tables cleaned.")
}
