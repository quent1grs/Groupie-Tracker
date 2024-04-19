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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des données du formulaire", http.StatusInternalServerError)
		return
	}
	username := r.FormValue("signname")
	email := r.FormValue("signemail")
	password := r.FormValue("signpass")
	// if password != confirmpassword {
	// 	http.Error(w, "Les mots de passe ne correspondent pas", http.StatusBadRequest)
	// 	return
	// } else if len(password) < 12 {
	// 	http.Error(w, "Le mot de passe doit contenir au moins 12 caractères", http.StatusBadRequest)
	// 	return
	// } else if len(username) < 3 {
	// 	http.Error(w, "Le nom d'utilisateur doit contenir au moins 3 caractères", http.StatusBadRequest)
	// 	return
	// }
	password = database.Hash(password)

	// Insérer les données dans la base de données en appelant la fonction existante
	err = database.InsertFormData(email, username, password)
	if err != nil {
		http.Error(w, "Erreur lors de l'insertion des données dans la base de données", http.StatusInternalServerError)
		return
	}
	//Actualiser la page en renvoyant le même fichier HTML
	http.ServeFile(w, r, "./home-page.html")
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return

	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des données du formulaire", http.StatusInternalServerError)
		return
	}

	// emailorUsername := r.FormValue("logemail/logname")
	password := r.FormValue("logpass")

	password = database.Hash(password)

	http.ServeFile(w, r, "./choosegamepage.html")
}

func UpdateUser(db *sql.DB, id int, username string, email string, password string) {
	stmt, err := db.Prepare("UPDATE USER SET username = ?, email = ?, password = ? WHERE id = ?")

	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(username, email, password, id)

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
