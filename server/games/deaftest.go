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
	// defer conn.Close()

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
			var wsResponse map[string]interface{}
			err := json.Unmarshal(message, &wsResponse)
			if err != nil {
				log.Println("Error parsing message:", err)
				return
			}

			userResponse, resp := wsResponse["answer"].(string)
			remainingTime, time := wsResponse["remainingTime"].(float64)
			if !resp || !time {
				println("Error getting data from message:", wsResponse)
			}

			status, response := CheckAnswer(userResponse, currentMusic)
			if status {
				user.CorrectAnswer = true
				if userResponse == currentMusic.Title || userResponse == currentMusic.Artist {
					user.Score += int(remainingTime)
				} else if userResponse == currentMusic.Artist+" "+currentMusic.Title || userResponse == currentMusic.Title+" "+currentMusic.Artist {
					user.Score += int(remainingTime) + 5
				}
			}

			SendMessage(response, conn)
		}

		if string(message) == "Change_song" {
			NextMusic(&currentMusic, &music)
		}
	}
}
