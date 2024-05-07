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

type UserTable struct {
	Users map[string]*User
}

type User struct {
	Token         string
	Score         int
	CorrectAnswer bool
}

// open a websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients []*websocket.Conn
var currentMusic PageData

// load the blindtest page
func HandleBlindtest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/blindtest.html")
}

// function to handle the game logic
func BlindtestWs(w http.ResponseWriter, r *http.Request) {
	UserTable := NewUserTable()
	var music spotifyapi.Music

	if currentMusic.Artist == "" {
		body := spotifyapi.GetPlaylist("https://api.spotify.com/v1/playlists/3hhUZQwNteEDClZTu4XY9X")
		music = spotifyapi.ParsePlaylist(body)
		currentMusic = ChooseMusic(music)
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

		println(string(message))
		if !user.CorrectAnswer {
			if string(message) != "start_music" && string(message) != "end_music" {
				userResponse := string(message)
				var response map[string]string
				if userResponse == currentMusic.Artist || userResponse == currentMusic.Title || userResponse == currentMusic.Artist+currentMusic.Title || userResponse == currentMusic.Title+currentMusic.Artist {
					response = map[string]string{
						"message": "Correct!",
					}
					user.Score += 1
					user.CorrectAnswer = true
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
			NextMusic(&currentMusic, &music)
		}
	}
}

func NewUserTable() *UserTable {
	return &UserTable{
		Users: make(map[string]*User),
	}
}

func (userTable *UserTable) AddUser(cookie string, user *User) {
	userTable.Users[cookie] = user
}

func (ut *UserTable) GetUser(cookie string) (*User, bool) {
	user, exists := ut.Users[cookie]
	return user, exists
}

func ChooseMusic(music spotifyapi.Music) PageData {
	i := mrand.Intn(len(music.Artists))
	PageData := PageData{
		Artist:  music.Artists[i],
		Title:   music.Titles[i],
		Preview: music.MusicPreview[i],
		Lyrics:  music.MusicLyrics[i],
	}

	music.Artists = append(music.Artists[:i], music.Artists[i+1:]...)
	music.Titles = append(music.Titles[:i], music.Titles[i+1:]...)
	music.MusicPreview = append(music.MusicPreview[:i], music.MusicPreview[i+1:]...)
	music.MusicLyrics = append(music.MusicLyrics[:i], music.MusicLyrics[i+1:]...)

	return PageData
}

func NextMusic(currentMusic *PageData, music *spotifyapi.Music) {
	if len(music.Artists) == 0 {
		log.Println("No more music to play")
		return
	}

	*currentMusic = ChooseMusic(*music)

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
}
