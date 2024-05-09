package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
)

func Database() {
	// user_id := 0
	db, err := sql.Open("sqlite3", "database/database.go")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
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

// InsertFormData insère les données du formulaire dans la base de données
func InsertFormData(username, email, password string) error {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture de la base de données: %v", err)
	}

	stmt, err := db.Prepare("INSERT INTO USER(id, username, email, password, status, sessioncookie) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("erreur lors de la préparation de la requête d'insertion: %v", err)
	}
	defer stmt.Close()

	// Note : keep email and username in this order, even if it's in reverse order in the database
	_, err = stmt.Exec(nextAvailableID(), email, username, password, "offline", "")
	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion dans la base de données: %v", err)
	}

	fmt.Println("Données du formulaire insérées avec succès dans la base de données.")
	db.Close()

	return nil

}

func Hash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func GetUsers(db *sql.DB) {
	// Préparation de la requête SQL
	rows, err := db.Query("SELECT id, username, email, password FROM USER")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Parcours des lignes de résultat
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
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		return 0
	}
	defer db.Close()

	var maxID int
	err = db.QueryRow("SELECT MAX(id) FROM USER").Scan(&maxID)
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

func IsUsernameInDB(username string) bool {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT username FROM USER WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return rows.Next()
}

func IsEmailInDB(email string) bool {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	var hashedPassword string
	var idOfRow int
	err = db.QueryRow("SELECT password FROM USER WHERE username = ? OR email = ?", identifier, identifier).Scan(&hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow("SELECT id FROM USER WHERE username = ? OR email = ?", identifier, identifier).Scan(&idOfRow)
	if err != nil {
		log.Fatal(err)
	}

	// Si identifier ne correspond pas à l'username ou à l'email, on renvoit false
	if !isIdentifierPresentInTheRow(identifier, hashedPassword, idOfRow) {
		return false
	}

	var comparedPassword string
	err = db.QueryRow("SELECT password FROM USER WHERE id = ?", idOfRow).Scan(&comparedPassword)
	if err != nil {
		log.Fatal(err)
	}

	return Hash(password) == comparedPassword
}

func isIdentifierPresentInTheRow(identifier string, hashedPassword string, id int) bool {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	var comparedUsername string
	var comparedEmail string
	err = db.QueryRow("SELECT username FROM USER WHERE id = ?", id).Scan(&comparedUsername)
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
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE USER SET status = 'offline'")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("UPDATE USER SET sessioncookie = ''")
	if err != nil {
		log.Fatal(err)
	}
}
