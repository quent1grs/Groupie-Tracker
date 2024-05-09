package games

import "net/http"

// Manages the GAMES table from db.sqlite database.

// Game struct

type Game struct {
	id       string
	name     string
	gameName string
}

// func HandleCreateGame(w http.ResponseWriter, r *http.Request) {
// 	// Récupérer le nom de l'utilisateur par le cookie
// 	username := session.GetUsernameFromCookie(r)
// 	// Récupérer les données du formulaire
// 	r.ParseForm()
// }

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
