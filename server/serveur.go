package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"groupietracker/database"
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

	"github.com/fr05t1k/musixmatch/entity/lyrics"
	_ "github.com/mattn/go-sqlite3"
)

// DEF CONFIGURABLES
var PORT = "8080"
var HOST = ""
var addr = flag.String("addr", HOST+":"+PORT, "http service address")

// ENDEF CONFIGURABLES

type PageData struct {
	URL    string
	Lyrics *lyrics.Lyrics
}

// var loggedUsers = make(map[string]user.User)

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

	// title := playlist.Tracks.Items[i].Track.Name
	// artist := playlist.Tracks.Items[i].Track.Artists[0].Name

	// lyrics := spotifyapi.GetLyrics(title, artist)
	// fmt.Println(lyrics.Language)
	// data: variable à passer à la page HTML pour la musique
	// data := PageData{
	// 	URL:    musicUrl[i],
	// 	Lyrics: lyrics,
	// }

	fmt.Println("Launching server.")
	fmt.Println("Current server address: " + *addr)
	fs := http.FileServer(http.Dir("./assets"))
	if fs == nil {
		log.Fatal("File server not found.")
	} else {
		log.Printf("File server found.")
	}

	// Code issu de la démo de chat. À conserver pour le chat global.
	// hub := newHub()
	// go hub.run()
	// http.HandleFunc("/ws", func(w http.ResponseWriter, rhttp.Request) {
	//     serveWs(hub, w, r)
	// })

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/signup", user.HandleSignup)
	http.HandleFunc("/login", user.HandleLogin)
	http.HandleFunc("/", handleHome)
	// TODO : Routes à ajouter

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
