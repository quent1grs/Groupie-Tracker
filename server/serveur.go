package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"groupietracker/database"
	"groupietracker/server/lobby"
	session "groupietracker/server/session"
	"io"
	"log"
	mrand "math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	user "groupietracker/server/user"
	spotifyapi "groupietracker/spotifyApi"

	_ "github.com/mattn/go-sqlite3"
)

// DEF CONFIGURABLES
var PORT = "8080"
var HOST = ""
var addr = flag.String("addr", HOST+":"+PORT, "http service address")

// ACTIVE OBJECTS
var loggedUsers = make(map[string]string)
var activeCookies = make(map[string]session.Cookie)
var activeSessions = make(map[string]session.Session)

type PageData struct {
	URL string
}

func main() {
	musicUrl := []string{}

	token := getToken()
	body := spotifyapi.GetPlaylist("https://api.spotify.com/v1/playlists/3hhUZQwNteEDClZTu4XY9X", token)

	var playlist spotifyapi.SearchResponse
	err := json.Unmarshal(body, &playlist)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range playlist.Tracks.Items {
		uri := path.Base(item.Track.ExternalUrls.Spotify)
		musicUrl = append(musicUrl, uri)
	}
	i := mrand.Intn(len(musicUrl))
	println("url de la musique : " + musicUrl[i])

	fmt.Println("Launching server.")
	fmt.Println("Current server address: " + *addr)
	fs := http.FileServer(http.Dir("./assets"))
	if fs == nil {
		log.Fatal("File server not found.")
	} else {
		log.Printf("File server found.")
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
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
			for _, session := range activeSessions {
				if time.Now().Unix()-session.InactiveSince > 360 {
					fmt.Println("Session inactive depuis 6 minutes. Suppression de la session.")
					fmt.Println("Session : ", session)
					delete(activeSessions, session.Username)
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
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Printf("Server listening at " + *addr)
	log.Fatal(server.ListenAndServe())

	database.Database()

}

func getToken() string {
	clientID := "c27ae1942ee94d23a21f324b6feba015"
	clientSecret := "c527485ba55545a4a0e88614a886500a" // Base64 encode the client ID and secret
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return result["access_token"].(string)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./home-page.html")
}

func IsCookieActive(cookie session.Cookie) bool {
	for _, c := range activeCookies {
		if c.CookieID == cookie.CookieID {
			return true
		}
	}
	return false
}
