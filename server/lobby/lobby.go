package lobby

import (
	"fmt"
	session "groupietracker/server/session"
	"net/http"
)

func HandleLobby(w http.ResponseWriter, r *http.Request) {
	if !session.IsCookieActive(session.Cookie{CookieID: r.Header.Get("Cookie")}) {
		fmt.Println("Cookie : ", r.Header.Get("Cookie"))
		fmt.Println("User not logged in. Redirecting to login page.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		fmt.Println("User logged in. Proceeding to lobby.")
	}

	fmt.Println("Cookie : ", r.Header.Get("Cookie"))
	fmt.Println("User : ", session.ActiveSessions[r.Header.Get("Cookie")].Username)
	fmt.Println("Active sessions : ", session.ActiveSessions)

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./pages/choosegamepage.html")
}
