package games

import (
	"fmt"
	gamesdb "groupietracker/database/gamesDB"
	roomusersdb "groupietracker/database/roomusersDB"
	userdb "groupietracker/database/userDB"
	"groupietracker/server/session"
	"net/http"
)

// Manages the GAMES table from db.sqlite database.

// Game struct

type Game struct {
	id       string
	name     string
	gameName string
}

func HandleGame(w http.ResponseWriter, r *http.Request) {
	// Identifier le jeu à partir de l'username
	fmt.Println("[DEBUG] games.HandleGame() called.")
	username := session.GetUsername(w, r)
	fmt.Println("[DEBUG] Username : " + username)

	userid := userdb.GetIDFromUsername(username)
	fmt.Println("[DEBUG] User ID : " + string(rune(userid)))

	room := roomusersdb.GetRoomAssociatedWithUser(userid)
	fmt.Println("[DEBUG] Room : " + string(rune(room)))

	gameMode := gamesdb.GetGameMode(room)
	fmt.Println("[DEBUG] Game mode : " + gameMode)

	if gameMode == "scattegories" {
		HandleScattegories(w, r)
	} else if gameMode == "deaftest" {
		HandleDeaftest(w, r)
	} else if gameMode == "blindtest" {
		HandleBlindtest(w, r)
	}

}

func HandleDeleteGame() {

}

// Problématique de la fonction : gestion de chaque room de jeu à l'ID unique selon le format "/room/{roomID}"
// Chaque room doit pouvoir être gérée de manière indépendante des autres rooms. Les informations nécessaires sont
// stockées dans la base de données (ROOM_USERS)
func HandleRoom(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de la room
	// Récupérer les informations de la room
	// Afficher les informations de la room
	// Gérer les actions de la room
}
