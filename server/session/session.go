package session

import (
	"database/sql"
	"fmt"
	"groupietracker/database"
	"log"
	"math/rand"
	"net/http"
	"strings"
	// "groupietracker/server"
)

// var LoggedUsers = make(map[string]string)
// var ActiveCookies = make(map[string]Cookie)
// var ActiveSessions = make(map[string]Session)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const cookieSize = 20

func IssueCookie() string {
	return generateCookieID()
}

func IsCookieValid(cookie string) bool {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT sessioncookie FROM USER")
	if err != nil {
		fmt.Println("[DEBUG] Error while querying database.")
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var c string
		err = rows.Scan(&c)
		if err != nil {
			log.Fatal(err)
		}
		if c == cookie {
			return true
		}
	}
	return false
}

func generateCookieID() string {
	b := make([]byte, cookieSize)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func IsClientLoggedIn(r *http.Request) bool {
	username := strings.Split(strings.Split(r.Header.Get("Cookie"), "; ")[0], "=")[1]
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT status FROM USER WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		err = rows.Scan(&status)
		if err != nil {
			log.Fatal(err)
		}
		if status == "connected" {
			return true
		}
	}
	return false
}

func SetStatusConnected(w http.ResponseWriter, r *http.Request) {
	username := strings.Split(strings.Split(r.Header.Get("Cookie"), "; ")[0], "=")[1]
	// cookie := strings.Split(strings.Split(r.Header.Get("Cookie"), "; ")[1], "=")[1]
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE USER SET status = 'connected' WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie := r.Header.Get("Cookie")

	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE USER SET status = 'disconnected' WHERE sessioncookie = ?", cookie)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM USER WHERE sessioncookie = ?", cookie)
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Del("Cookie")
	r.Header.Add("Cookie", "cookie=deleted")
	http.SetCookie(w, &http.Cookie{Name: "cookie", Value: "deleted", MaxAge: -1})
	http.ServeFile(w, r, "./home-page.html")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UpdateCookieInDB(cookie string, username string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE USER SET sessioncookie = ? WHERE username = ?", cookie, username)
	if err != nil {
		log.Fatal(err)
	}
	database.ShowUserDetails(username)
	fmt.Println("[DEBUG] New cookie: ", cookie)
	fmt.Println("[DEBUG] Cookie updated in database.")
}

func GetUsername(w http.ResponseWriter, r *http.Request) string {
	return strings.Split(strings.Split(r.Header.Get("Cookie"), "; ")[0], "=")[1]
}

func GetCookie(w http.ResponseWriter, r *http.Request) string {
	return strings.Split(strings.Split(r.Header.Get("Cookie"), "; ")[1], "=")[1]
}
