package main

import (
	"flag"
	"fmt"
	"groupietracker/database"
	"groupietracker/server/lobby"
	session "groupietracker/server/session"
	"log"
	"net/http"
	"time"

	"groupietracker/server/games"
	user "groupietracker/server/user"
	spotifyapi "groupietracker/spotifyApi"

	_ "github.com/mattn/go-sqlite3"
)

// DEF CONFIGURABLES
var PORT = "8080"
var HOST = ""
var addr = flag.String("addr", HOST+":"+PORT, "http service address")

type PageData struct {
	URL string
}

func main() {
	body := spotifyapi.GetPlaylist("https://api.spotify.com/v1/playlists/3hhUZQwNteEDClZTu4XY9X")

	spotifyapi.ParsePlaylist(body)

	fmt.Println("Launching server.")
	fmt.Println("Current server address: " + *addr)
	fs := http.FileServer(http.Dir("./assets"))
	if fs == nil {
		log.Fatal("File server not found.")
	} else {
		log.Printf("File server found.")
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/blindtest", games.HandleBlindtest)
	// http.HandleFunc("/deaftest", games.HandleDeafTest)
	http.HandleFunc("/scattegories", games.HandleScattegories)
	http.HandleFunc("/signup", user.HandleSignup)
	http.HandleFunc("/login", user.HandleLogin)
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/checkUsername", user.HandleCheckUsername) // Nouvelle route : checkUsername (pour vérifier la disponibilité d'un nom d'utilisateur lors de l'inscription par requête AJAX)
	http.HandleFunc("/checkEmail", user.HandleCheckEmail)       // Nouvelle route : checkEmail (pour vérifier la disponibilité d'un email lors de l'inscription par requête AJAX)
	http.HandleFunc("/isPasswordValid", user.HandleIsPasswordValid)
	http.HandleFunc("/isUserInDB", user.HandleIsUserInDB)
	http.HandleFunc("/isPasswordCorrect", user.HandleIsPasswordCorrect)
	http.HandleFunc("/lobby", lobby.HandleLobby)

	fmt.Println(time.Now().String() + " Server is running on port " + PORT)

	// Horloge pour les sessions
	go func() {
		for {
			time.Sleep(1 * time.Second)
			for _, session := range session.ActiveSessions {
				if time.Now().Unix()-session.InactiveSince > 360 {
					fmt.Println("Session inactive depuis 6 minutes. Suppression de la session.")
					fmt.Println("Session : ", session)
					// TODO
				} else {
					// Incrémenter la variable inactiveSince de chaque session de 1 seconde
					session.InactiveSince++
				}
			}
		}
	}()

	server := &http.Server{
		Addr:              *addr,
		Handler:           nil,
		ReadHeaderTimeout: 3 * time.Second,
	}

	// Démarrage du serveur
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Printf("Server listening at " + *addr)
	log.Fatal(server.ListenAndServe())

	database.Database()

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./home-page.html")
}

func IsCookieActive(cookie session.Cookie) bool {
	for _, c := range session.ActiveSessions {
		if c.Cookie.CookieID == cookie.CookieID {
			return true
		}
	}
	return false
}
