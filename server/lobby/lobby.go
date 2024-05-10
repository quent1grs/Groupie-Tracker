package lobby

import (
	session "groupietracker/server/session"
	"net/http"
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

	if r.Header.Get("cookie") == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cookie := session.GetCookie(w, r)

	if !session.IsCookieValid(cookie) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./pages/choosegamepage.html")
}
