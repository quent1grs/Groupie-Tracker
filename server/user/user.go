package user

import (
	"database/sql"
	"encoding/base64"
	"log"
)

type User struct {
	Id             int
	Username       string
	Password       string
	Email          string
	ProfilePicture string
	// PersonalPlaylist []Playlist
	Experience int
}

// Ajout d'un utilisateur à la base de données
func CreateUser(db *sql.DB, id int, pseudo string, password string, email string) {
	stmt, err := db.Prepare("INSERT INTO USER(id, pseudo, email, password) VALUES(?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id, pseudo, email, password)

	if err != nil {
		log.Fatal(err)
	}
}

func UpdateUser(db *sql.DB, id int, pseudo string, password string, email string) {
	stmt, err := db.Prepare("UPDATE USER SET pseudo = ?, email = ?, password = ? WHERE id = ?")

	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(pseudo, email, password, id)

	if err != nil {
		log.Fatal(err)
	}
}

func GetUser(db *sql.DB, id int) User {
	var user User
	stmt, err := db.Prepare("SELECT * FROM USER WHERE id = ?")

	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(id).Scan(&user.Username, &user.Email, &user.Password)

	if err != nil {
		log.Fatal(err)
	}
	return user
}

func RemoveUser(db *sql.DB, id int) {
	stmt, err := db.Prepare("DELETE FROM USER WHERE id = ?")

	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id)

	if err != nil {
		log.Fatal(err)
	}
}

func EncodePassword(password string) string {
	// Encode password
	password = base64.StdEncoding.EncodeToString([]byte(password))
	return password
}

func DecodePassword(password string) string {
	// Decode password
	decoded, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		log.Fatal(err)
	}
	return string(decoded)
}
