package ws

import (
	"fmt"
	"net/http"

	session "groupietracker/server/session"

	"github.com/gorilla/websocket"
)

// Manages the WS logic for the different rooms/games.

// Stores the clients in a room.
type Room struct {
	usernames map[string]bool
	clients   map[*websocket.Conn]bool
}

// Stores the rooms.
var Rooms = make(map[string]Room)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// Select the websocket connection and the room ID.
	// roomID := session.GetRoomIDCookie(w, r)
	// conn := GetConnection(w, r)

	// OR

	// vars := mux.Vars(r)
	// roomID := vars["roomID"]
}

func CreateNewWSAndAddToRoom(w http.ResponseWriter, r *http.Request, roomID string) {
	fmt.Println("[DEBUG] ws.CreateNewWSAndAddToRoom() called.")

	conn, _ := upgrader.Upgrade(w, r, nil)

	Rooms[roomID] = Room{
		clients: make(map[*websocket.Conn]bool),
	}

	Rooms[roomID].clients[conn] = true
}

func AddNewRoom(roomID string) {
	fmt.Println("[DEBUG] ws.AddNewRoom() called.")

	var newRoom Room
	newRoom.clients = make(map[*websocket.Conn]bool)
	Rooms[roomID] = newRoom

	showRooms()
}

func GetConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	// Récupère le nom de l'utilisateur
	username := session.GetUsername(w, r)
	// Récupère le roomID
	roomID := session.GetRoomIDCookie(w, r)

	for c := range Rooms[roomID].clients {
		if Rooms[roomID].usernames[username] {
			return c
		}
	}

	return nil
}

// DEBUG FUNCTIONS
func showRooms() {
	for k, v := range Rooms {
		println("Room ID : " + k)
		for c := range v.clients {
			println("Client : " + c.LocalAddr().String())
		}
	}
}
