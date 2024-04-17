package user

import (
	"database/sql"
	"encoding/base64"
	"groupietracker/database"
	"log"
	"net/http"
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
func HandleSignup(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return

	}
	// Récupérer les données du formulaire depuis la requête HTTP
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des données du formulaire", http.StatusInternalServerError)
		return
	}
	email := r.FormValue("signname")
	password := r.FormValue("signemail")
	username := r.FormValue("signpass")

	// Insérer les données dans la base de données en appelant la fonction existante
	err = database.InsertFormData(email, password, username)
	if err != nil {
		http.Error(w, "Erreur lors de l'insertion des données dans la base de données", http.StatusInternalServerError)
		return
	}
	//Actualiser la page en renvoyant le même fichier HTML
	http.ServeFile(w, r, "./home-page.html")
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
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
