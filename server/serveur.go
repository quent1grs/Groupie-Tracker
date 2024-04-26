package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"groupietracker/database"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
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
	var Music spotifyapi.Music
	token := getToken()
	body := spotifyapi.GetPlaylist("https://api.spotify.com/v1/playlists/3hhUZQwNteEDClZTu4XY9X", token)

	Music.Artists, Music.Titles, Music.MusicUrl, Music.MusicLyrics = parsePlaylist(body)

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

func parsePlaylist(body []byte) ([]string, []string, []string, []string) {
	var musicUrl []string
	var musicLyrics []string
	var artists []string
	var titles []string
	var playlist spotifyapi.SearchResponse

	err := json.Unmarshal(body, &playlist)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range playlist.Tracks.Items {
		uri := path.Base(item.Track.ExternalUrls.Spotify)
		musicUrl = append(musicUrl, uri)
		title := item.Track.Name
		titles = append(titles, title)
		artist := item.Track.Artists[0].Name
		artists = append(artists, artist)
		lyrics := spotifyapi.GetLyrics(title, artist)
		parts := strings.Split(lyrics.Body, "\n...\n\n******* This Lyrics is NOT for Commercial use *******")
		lyrics.Body = parts[0]
		musicLyrics = append(musicLyrics, lyrics.Body)
	}

	return artists, titles, musicUrl, musicLyrics
}
