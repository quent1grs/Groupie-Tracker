package session

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	// "groupietracker/server"
)

// var LoggedUsers = make(map[string]string)
// var ActiveCookies = make(map[string]Cookie)
var ActiveSessions = make(map[string]Session)

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

func IssueCookie(username string) Cookie {
	return Cookie{CookieToken: generateCookieID(), Username: username}
}

func IsCookieActive(cookie Cookie) bool {
	// Check if the cookie is in the activeCookies map
	// If it is, return true
	// If it isn't, return false
	fmt.Println("[DEBUG] Cookie to check: ", cookie)
	for _, c := range ActiveSessions {
		// fmt.Println("Current from-list Cookie ID: ", c.Cookie.CookieID)
		// fmt.Println("Examined Cookie ID: ", cookie.CookieID)
		if "cookie="+c.Cookie.CookieToken == cookie.CookieToken {
			fmt.Println("[DEBUG] Cookie is active.")
			return true
		}
	}
	fmt.Println("[DEBUG] Cookie is not active.")
	return false
}

func generateCookieID() string {
	b := make([]byte, cookieSize)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	fmt.Println("Generated cookie ID: ", string(b))
	return string(b)
}

func IsClientLoggedIn(r *http.Request) bool {
	// Check if the client is logged in
	// If the client is logged in, return true
	// If the client is not logged in, return false
	cookie := r.Header.Get("Cookie")
	fmt.Println("[DEBUG] Cookie: ", cookie)

	yesNo := IsCookieActive(Cookie{CookieToken: cookie})
	if yesNo == true {
		fmt.Println("[DEBUG] Client is logged in.")
	} else {
		fmt.Println("[DEBUG] Client is not logged in.")
	}
	return yesNo
}

func AddSession(username string, cookie Cookie) {
	// Add the session to the ActiveSessions map
	// The key should be the username
	// The value should be the session
	// The session should have the current time as the InactiveSince field
	// Si la map ActiveSessions n'existe pas, la créer

	ActiveSessions[username] = Session{Username: username, Cookie: cookie, InactiveSince: 0}
	// énumération des sessions actives
	for key, value := range ActiveSessions {
		fmt.Println("Key:", key, "Value:", value)
	}

}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Handling logout.")
	cookie := r.Header.Get("Cookie")
	fmt.Println("[DEBUG] Cookie: ", cookie)

	// Retire le cookie sur la session locale du client

	fmt.Println("[DEBUG] Active sessions: ", ActiveSessions)
	for key, value := range ActiveSessions {
		fmt.Println("[DEBUG] Key:", key, "Value:", value)
		if "cookie="+value.Cookie.CookieToken == cookie {
			delete(ActiveSessions, key)
		}
	}
	fmt.Println("[DEBUG] Active sessions after removal of cookie : ", ActiveSessions)
	r.Header.Del("Cookie")
	r.Header.Add("Cookie", "cookie=deleted")
	fmt.Println("[DEBUG] Cookie after logout: ", r.Header.Get("Cookie"))
	http.SetCookie(w, &http.Cookie{Name: "cookie", Value: "deleted", MaxAge: -1})

	http.ServeFile(w, r, "./home-page.html")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func GetUsername(w http.ResponseWriter, r *http.Request) string {
	return strings.Split(strings.Split(r.Header.Get("Cookie"), "; ")[0], "=")[1]
}

func GetCookie(w http.ResponseWriter, r *http.Request) string {
	return strings.Split(strings.Split(r.Header.Get("Cookie"), "; ")[1], "=")[1]
}
