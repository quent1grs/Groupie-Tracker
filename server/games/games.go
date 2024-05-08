package games

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
