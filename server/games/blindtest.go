package games

import (
	"encoding/json"
	spotifyapi "groupietracker/spotifyApi"
	"log"
	mrand "math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

type PageData struct {
	Artist  string
	Title   string
	Preview string
	Lyrics  string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients []*websocket.Conn
var currentMusic PageData

func HandleBlindtest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/blindtest.html")
}

func BlindtestWs(w http.ResponseWriter, r *http.Request) {
	if currentMusic.Artist == "" {
		body := spotifyapi.GetPlaylist("https://api.spotify.com/v1/playlists/3hhUZQwNteEDClZTu4XY9X")
		music := spotifyapi.ParsePlaylist(body)
		i := mrand.Intn(len(music.Artists))
		currentMusic = PageData{
			Artist:  music.Artists[i],
			Title:   music.Titles[i],
			Preview: music.MusicPreview[i],
			Lyrics:  music.MusicLyrics[i],
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	clients = append(clients, conn)

	jsonData, err := json.Marshal(currentMusic)
	if err != nil {
		log.Println(err)
		return
	}

	for _, client := range clients {
		err = client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)

			// Supprimez le client de la liste des clients
			for i, client := range clients {
				if client == conn {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}

			break
		}

		if string(message) == "start_music" {
			for _, client := range clients {
				err = client.WriteMessage(websocket.TextMessage, []byte("start_music"))
				if err != nil {
					log.Println("write:", err)
					continue
				}
			}
		}

		if string(message) == "end_music" {
			for _, client := range clients {
				err = client.WriteMessage(websocket.TextMessage, []byte("end_music"))
				if err != nil {
					log.Println("write:", err)
					continue
				}
			}
		}

		if string(message) == "answer" {
			userResponse := string(message)
			var response map[string]string
			if userResponse == currentMusic.Artist || userResponse == currentMusic.Title || userResponse == currentMusic.Artist+currentMusic.Title || userResponse == currentMusic.Title+currentMusic.Artist {
				response = map[string]string{
					"message": "Correct!",
				}
			} else {
				response = map[string]string{
					"message": "Incorrect!",
				}
			}

			jsonData, err := json.Marshal(response)
			if err != nil {
				log.Println("Error marshalling response:", err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Println("Error writing message:", err)
				break
			}
		}

		log.Println("Received message:", string(message))

	}
}
