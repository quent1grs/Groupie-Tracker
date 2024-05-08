package lobby

import (
	"fmt"
	"groupietracker/database"
	session "groupietracker/server/session"
	"net/http"
	"strings"
)

func HandleLobby(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/lobby" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer la valeur du cookie dont le nom est "cookie"
	var cookieContent string
	if r.Header.Get("cookie") == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	cookieContent = r.Header.Get("cookie")

	// username := strings.Split(strings.Split(cookieContent, "; ")[0], "=")[1]
	cookie := strings.Split(strings.Split(cookieContent, "; ")[1], "=")[1]

	if !session.IsCookieValid(cookie) {
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

	fmt.Println("Active sessions : ")
	database.EnumerateConnectedUsers()

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "pages/choosegamepage.html")
}
