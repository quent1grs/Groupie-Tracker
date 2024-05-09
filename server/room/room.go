package room

// Stores logic for rooms management. A room is a group session with a dedicated url, it's own chat and users list.
// Rooms are stored in the database in ./database/db.sqlite in the table 'rooms'.

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	// "groupietracker/server/room/roomUsers"
	gamesdb "groupietracker/database/gamesDB"
	roomsdb "groupietracker/database/roomsDB"
	roomusersdb "groupietracker/database/roomusersDB"
	userdb "groupietracker/database/userDB"
	"groupietracker/server/room/roomUsers"
	"groupietracker/server/session"
	"log"
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
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM ROOMS")
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteRoomFromDB(roomID string) {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM ROOMS WHERE id = ?", roomID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[EVENT] Room" + roomID + " deleted from database.")
}

func doesRoomExist(roomID string) bool {
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id FROM ROOMS WHERE id = ?", roomID)
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
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT COUNT(*) FROM ROOMS_USERS WHERE roomID = ?", roomID)
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
