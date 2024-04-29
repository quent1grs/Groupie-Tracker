package session

import (
	"fmt"
	"math/rand"
	// "groupietracker/server"
)

// var LoggedUsers = make(map[string]string)
// var ActiveCookies = make(map[string]Cookie)
var ActiveSessions = make(map[string]Session)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const cookieSize = 20

type Cookie struct {
	CookieID string
}

type Session struct {
	Username      string
	Cookie        Cookie
	InactiveSince int64
}

func IssueCookie(username string) Cookie {
	return Cookie{CookieID: generateCookieID()}
}

func IsCookieActive(cookie Cookie) bool {
	// Check if the cookie is in the activeCookies map
	// If it is, return true
	// If it isn't, return false
	for _, c := range ActiveSessions {
		// fmt.Println("Current from-list Cookie ID: ", c.Cookie.CookieID)
		// fmt.Println("Examined Cookie ID: ", cookie.CookieID)
		if "cookie="+c.Cookie.CookieID == cookie.CookieID {
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
	fmt.Println("Generated cookie ID: ", string(b))
	return string(b)
}

// func IsClientLoggedIn(cookie Cookie) bool {
// 	// Check if the cookie is in the activeCookies map
// 	// If it is, return true
// 	// If it isn't, return false
// 	for _, c := range activeCookies {
// 		if c.CookieID == cookie.CookieID {
// 			return true
// 		}
// 	}
// 	return false
// }

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
