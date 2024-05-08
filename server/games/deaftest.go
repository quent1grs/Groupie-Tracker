package games

import (
	"encoding/json"
	spotifyapi "groupietracker/spotifyApi"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleDeaftest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/deaftest.html")
}

func DeaftestWs(w http.ResponseWriter, r *http.Request) {
	UserTable := NewUserTable()
	var music spotifyapi.Music

	if currentMusic.Artist == "" {
		body := spotifyapi.GetPlaylist("https://api.spotify.com/v1/playlists/3hhUZQwNteEDClZTu4XY9X")
		music = spotifyapi.ParsePlaylist(body)
		currentMusic = ChooseMusic(&music)
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
		println(string(message))
		if err != nil {
			log.Println("Error reading message:", err)

			for i, client := range clients {
				if client == conn {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			break
		}

		user, exist := UserTable.GetUser(r.Header.Get("Cookie"))
		if !exist {
			user = &User{
				Token: r.Header.Get("Cookie"),
				Score: 0,
			}
			UserTable.AddUser(r.Header.Get("Cookie"), user)
		}

		if !user.CorrectAnswer {
			userResponse := string(message)
			var response map[string]string
			status, response := CheckAnswer(userResponse, currentMusic)
			if status {
				user.CorrectAnswer = true
				user.Score += 1
			}

			SendMessage(response, conn)
		}

		if string(message) == "Change_song" {
			NextMusic(&currentMusic, &music)
		}
	}
}
