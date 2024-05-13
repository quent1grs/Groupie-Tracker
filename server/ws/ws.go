package ws

import (
	"encoding/json"
	"log"
	"net/http"

	session "groupietracker/server/session"

	"github.com/gorilla/websocket"
)

// Manages the WS logic for the different rooms/games.

// Stores the clients in a room.
type Room struct {
	Usernames map[string]bool
	Clients   map[*websocket.Conn]bool
	ClientsWs []*websocket.Conn
}

// Stores the rooms.
var Rooms = make(map[string]Room)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Handles the WS connection sent by the client thru the "/ws" endpoint, reached by javascript scripts.
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erreur lors de l'upgrade de la connexion WebSocket:", err)
		return
	}
	defer conn.Close()

	// Récupère le nom de l'utilisateur et le roomID
	username := session.GetUsername(w, r)
	roomID := session.GetRoomIDCookie(w, r)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil { // Si l'utilisateur se déconnecte,
			delete(Rooms[roomID].Clients, conn)
			log.Printf("Déconnexion de l'utilisateur %s de la room %s", username, roomID)
			break
		}

		// Analyser le message JSON
		var message struct {
			Answer        string `json:"answer"`
			RemainingTime int    `json:"remainingTime"`
			Game          string `json:"game"`
		}
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Erreur lors de l'analyse du message JSON:", err)
			continue
		}

		// Diriger le message vers le gestionnaire approprié en fonction du jeu
		switch message.Game {
		case "deaftest":
			http.Redirect(w, r, "/deaftestws", http.StatusSeeOther)
		case "blindtest":
			http.Redirect(w, r, "/blindtestws", http.StatusSeeOther)
		case "scattegories":
			http.Redirect(w, r, "/ScattegoriesGame", http.StatusSeeOther)
		default:
			log.Printf("Jeu inconnu: %s", message.Game)
		}
	}
}

func CreateNewWSAndAddToRoom(w http.ResponseWriter, r *http.Request, roomID string) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	Rooms[roomID] = Room{
		Clients: make(map[*websocket.Conn]bool),
	}

	Rooms[roomID].Clients[conn] = true
}

func AddNewRoom(roomID string) {
	var newRoom Room
	newRoom.Clients = make(map[*websocket.Conn]bool)
	Rooms[roomID] = newRoom

}

// func GetConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Erreur lors de l'upgrade de la connexion WebSocket:", err)
// 	}
// 	roomID := session.GetRoomIDCookie(w, r)

// 	for c, _ := range Rooms[roomID].Clients {
// 		if c.LocalAddr().String() == conn.LocalAddr().String() {
// 			return c
// 		}
// 	} GroupieTracker0!

// 	return nil
// }

func GetRoom(roomID string) Room {
	return Rooms[roomID]
}

func (r *Room) GetConnections() map[*websocket.Conn]bool {
	return r.Clients
}

func (r *Room) GetClientsWs() []*websocket.Conn {
	return r.ClientsWs
}

// DEBUG FUNCTIONS
// func ShowRooms() {
// 	for k, v := range Rooms {
// 		println("Room ID : " + k)
// 		for c := range v.Clients {
// 			println("Client : " + c.LocalAddr().String())
// 		}
// 	}
// }
