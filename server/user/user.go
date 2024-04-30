package user

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"groupietracker/database"
	session "groupietracker/server/session"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"
)

type RequestNameRegisteringBody struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Password2 string `json:"password2"`
}

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
func HandleRegister(w http.ResponseWriter, r *http.Request) {

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

	password = database.Hash(password)

	if database.IsUsernameInDB(username) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if database.IsEmailInDB(email) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if PasswordEntropy(password) < 60 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Insérer les données dans la base de données en appelant la fonction existante
	err = database.InsertFormData(email, username, password)
	if err != nil {
		http.Error(w, "Erreur lors de l'insertion des données dans la base de données", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	emailorUsername := r.FormValue("logemail/loguser")
	password := r.FormValue("logpass")

	// // Si les informations de connexion ne sont pas correctes, rediriger vers la page de connexion
	if !database.IsPasswordCorrect(emailorUsername, password) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Si l'utilisateur a pu s'authentifier, créer un cookie et une session
	cookie := session.IssueCookie(emailorUsername)
	// Envoyer le cookie au client
	http.SetCookie(w, &http.Cookie{
		Name:  "cookie",
		Value: cookie.CookieID,
	})

	session.AddSession(emailorUsername, cookie)

	http.Redirect(w, r, "/lobby", http.StatusSeeOther)
}

func HandleLoginControl(w http.ResponseWriter, r *http.Request) {
	var requestNameRegisteringBody RequestNameRegisteringBody
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&requestNameRegisteringBody)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
		return
	}
	var identifier string
	if requestNameRegisteringBody.Email != "" {
		identifier = requestNameRegisteringBody.Email
	} else {
		identifier = requestNameRegisteringBody.Username
	}
	if database.IsPasswordCorrect(identifier, requestNameRegisteringBody.Password) {
		fmt.Fprint(w, true)
	} else {
		fmt.Fprint(w, false)
	}
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

func HandleRegisterControl(w http.ResponseWriter, r *http.Request) {
	// Fusionne l'ensemble des processus du contrôle d'inscription
	var requestNameRegisteringBody RequestNameRegisteringBody
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&requestNameRegisteringBody)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
		return
	}
	username := requestNameRegisteringBody.Username
	email := requestNameRegisteringBody.Email
	password := requestNameRegisteringBody.Password
	password2 := requestNameRegisteringBody.Password2
	matched, _ := regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", email)
	switch {
	case database.IsUsernameInDB(username):
		fmt.Fprint(w, "usernameInDB")
	case username == "":
		fmt.Fprint(w, "userEmpty")
	case database.IsEmailInDB(email):
		fmt.Fprint(w, "emailInDB")
	case !matched:
		fmt.Fprint(w, "emailInvalid")
	case email == "":
		fmt.Fprint(w, "emailEmpty")
	case PasswordEntropy(password) < 60:
		fmt.Fprint(w, "passwordEntropy")
	case password == "":
		fmt.Fprint(w, "passwordEmpty")
	case password != password2:
		fmt.Fprint(w, "passwordMismatchError")
	default:
		fmt.Fprint(w, true)
	}
}

// Vérifier si le nom d'utilisateur est disponible dans la base de données. Cette fonction est appelée par une requête AJAX
func HandleCheckUsername(w http.ResponseWriter, r *http.Request) {
	var requestNameRegisteringBody RequestNameRegisteringBody
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des données du formulaire", http.StatusInternalServerError)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&requestNameRegisteringBody)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
		return
	}
	username := requestNameRegisteringBody.Username

	if !database.IsUsernameInDB(username) {
		fmt.Fprint(w, false)
	} else {
		fmt.Fprint(w, true)
	}
}

// Vérifier si l'email est disponible dans la base de données
func HandleCheckEmail(w http.ResponseWriter, r *http.Request) {
	var requestNameRegisteringBody RequestNameRegisteringBody
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des données du formulaire", http.StatusInternalServerError)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&requestNameRegisteringBody)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
		return
	}
	email := requestNameRegisteringBody.Email
	matched, _ := regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", email)

	if !matched {
		fmt.Fprint(w, "invalid")
	} else if !database.IsEmailInDB(email) {
		fmt.Fprint(w, false)
	} else {
		fmt.Fprint(w, true)
	}
}

func PasswordEntropy(password string) int {
	// Calculate password entropy
	L := len(password)
	R := getCharsetPool(password)
	return int(math.Log2(math.Pow(float64(R), float64(L))))
}

func getCharsetPool(password string) int {
	// Get the charset pool of the password
	pool := 0
	if strings.ContainsAny(password, "0123456789") {
		pool += 10
	}
	if strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		pool += 26
	}
	if strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		pool += 26
	}
	if strings.ContainsAny(password, "`~!@#$%^&*()-=_+[{]}\\") {
		pool += 32
	}
	return pool
}

func HandleIsPasswordValid(w http.ResponseWriter, r *http.Request) {
	var requestNameRegisteringBody RequestNameRegisteringBody
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&requestNameRegisteringBody)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
		return
	}
	password := requestNameRegisteringBody.Password
	if PasswordEntropy(password) >= 60 {
		fmt.Fprint(w, true)
	} else {
		fmt.Fprint(w, false)
	}
}

func HandleIsUserInDB(w http.ResponseWriter, r *http.Request) {
	var requestNameRegisteringBody RequestNameRegisteringBody
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&requestNameRegisteringBody)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
		return
	}
	username := requestNameRegisteringBody.Username

	if database.IsUsernameInDB(username) {
		fmt.Fprint(w, "username")
	} else if database.IsEmailInDB(username) {
		fmt.Fprint(w, "email")
	} else {
		fmt.Fprint(w, false)
	}
}

func HandleIsPasswordCorrect(w http.ResponseWriter, r *http.Request) {
	var requestNameRegisteringBody RequestNameRegisteringBody
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&requestNameRegisteringBody)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
		return
	}
	username := requestNameRegisteringBody.Username
	password := requestNameRegisteringBody.Password

	if database.IsPasswordCorrect(username, password) {
		fmt.Fprint(w, true)
	} else {
		fmt.Fprint(w, false)
	}
}
