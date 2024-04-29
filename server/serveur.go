package main

import (
	"flag"
	"fmt"
	"groupietracker/database"
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

// ENDEF CONFIGURABLES

// var loggedUsers = make(map[string]user.User)

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
	// TODO : Routes à ajouter

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

	fmt.Println("Server launched.")
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./home-page.html")
}
