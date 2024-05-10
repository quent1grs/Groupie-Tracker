package session

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	db "groupietracker/database/DB_Connection" // Importation du package pour la connexion à la base de données
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const cookieSize = 20

type Cookie struct {
	Username    string
	CookieToken string
}

type Session struct {
	Username      string
	Cookie        Cookie
	InactiveSince int64
}

func IssueCookie() string {
	return generateCookieID()
}

func IsCookieValid(cookie string) bool {
	fmt.Println("[DEBUG] session.IsCookieValid() called.")
	fmt.Println("[DEBUG] Cookie : " + cookie)
	conn := db.GetDB()
	rows, err := conn.Query("SELECT sessioncookie FROM USER")
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
	fmt.Println("[DEBUG] Cookie is not valid.")
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
	cookie := r.Header.Get("Cookie")

	yesNo := IsCookieValid(cookie)
	if yesNo {
		return true
	} else {
		return false
	}
}

func UpdateCookieInDB(cookie string, username string) {
	conn := db.GetDB()
	_, err := conn.Exec("UPDATE USER SET sessioncookie = ? WHERE username = ?", cookie, username)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUsername(w http.ResponseWriter, r *http.Request) string {
	return strings.Split(strings.Split(r.Header.Get("Cookie"), "username=")[1], ";")[0]
}

func GetCookie(w http.ResponseWriter, r *http.Request) string {
	return strings.Split(strings.Split(r.Header.Get("Cookie"), "cookie=")[1], ";")[0]
}

func GetRoomIDCookie(w http.ResponseWriter, r *http.Request) string {
	return strings.Split(strings.Split(r.Header.Get("Cookie"), "roomIDCookie=")[1], ";")[0]
}
