package room

import (
	"encoding/json"
	"fmt"
	db "groupietracker/database/DB_Connection"
	gamesdb "groupietracker/database/gamesDB"
	roomsdb "groupietracker/database/roomsDB"
	roomusersdb "groupietracker/database/roomusersDB"
	userdb "groupietracker/database/userDB"
	"groupietracker/server/room/roomUsers"
	"groupietracker/server/session"
	"log"
	"net/http"
	"strconv"
)

type CreateRoomRequest struct {
	GameName      string   `json:"gameName"`
	GameType      string   `json:"gameType"`
	Rounds        string   `json:"rounds"`
	RoundDuration string   `json:"roundDuration"`
	MaxPlayers    string   `json:"maxPlayers"`
	Scattegories  []string `json:"scattegories"`
}

func HandleRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Path[6:]
	user := r.Header.Get("username")
	if !doesRoomExist(r.URL.Path[6:]) {
		fmt.Fprintln(w, "Room does not exist.")
	}
	if isRoomFull(roomID) {
		fmt.Fprintln(w, "Room is full.")
	}
	roomUsers.InsertUserInRoom(roomID, user)
	http.Redirect(w, r, "/room/"+roomID, http.StatusSeeOther)
}

func HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] Handling create room.")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	fmt.Printf("[DEBUG] Creating room...\n")
	user := session.GetUsername(w, r)
	var request CreateRoomRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	fmt.Println("[DEBUG] Request: ", request)
	fmt.Println(err)

	if err != nil {
		fmt.Println("[ERROR] Error while decoding request.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gameID := gamesdb.CreateGame(request.GameName)
	fmt.Println("[DEBUG] Game created with ID: ", gameID)
	roomID := roomsdb.GetNextAvailableID()
	fmt.Println("[DEBUG] Room created with ID: ", roomID)
	maxPlayers, _ := strconv.Atoi(request.MaxPlayers) // Convert string to int
	fmt.Println("[DEBUG] Max players: ", maxPlayers)
	roomsdb.InsertRoomInDatabase(roomID, user, maxPlayers, request.GameName, gameID)
	fmt.Println("[DEBUG] Room inserted in database.")

	roomusersdb.InsertUserInRoomUsers(roomID, userdb.GetIDFromUsername(user))
	fmt.Println("[DEBUG] User inserted in room users.")
	fmt.Fprint(w, "success="+string(rune(roomID)))
	http.Redirect(w, r, "/room/"+string(rune(roomID)), http.StatusSeeOther)
}

func CleaningRooms() {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection
	_, err := conn.Exec("DELETE FROM ROOMS")
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteRoomFromDB(roomID string) {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection
	_, err := conn.Exec("DELETE FROM ROOMS WHERE id = ?", roomID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] Room" + roomID + " deleted from database.")
}

func doesRoomExist(roomID string) bool {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection
	rows, err := conn.Query("SELECT id FROM ROOMS WHERE id = ?", roomID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}

func isRoomFull(roomID string) bool {
	conn := db.GetDB() // utilisez la fonction de connexion de votre package DB_Connection
	rows, err := conn.Query("SELECT COUNT(*) FROM ROOM_USERS WHERE id_room = ?", roomID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var count int
	rows.Scan(&count)
	if count >= 4 {
		return true
	}
	return false
}
