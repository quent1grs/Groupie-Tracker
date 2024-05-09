package main

import (
	"flag"
	"fmt"
	"groupietracker/database"
	gamesdb "groupietracker/database/gamesDB"
	roomsdb "groupietracker/database/roomsDB"
	roomusersdb "groupietracker/database/roomusersDB"
	"groupietracker/server/chat"
	"groupietracker/server/lobby"
	room "groupietracker/server/room"
	"log"
	"net/http"
	"time"

	"groupietracker/server/games"
	user "groupietracker/server/user"

<<<<<<< Updated upstream
	"github.com/gorilla/mux"
=======
	"github.com/gorilla/websocket"
>>>>>>> Stashed changes
	_ "github.com/mattn/go-sqlite3"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Ajoutez ceci si vous rencontrez des problèmes de CORS
}

// DEF CONFIGURABLES
var PORT = "8080"
var HOST = ""
var addr = flag.String("addr", HOST+":"+PORT, "http service address")

func main() {
	fmt.Println("Launching server...")
	fmt.Println("Initializing database...")
	database.ResetSessionData()
	gamesdb.ResetTable()
	roomsdb.DeleteAllRooms()
	roomusersdb.DeleteAllRoomUsers()
	chat.ChatScattegories()
	chat.ChatBlindtest()
	chat.ChatDeafTest()
	fmt.Println("Current server address: " + *addr)
	fs := http.FileServer(http.Dir("./assets"))
	if fs == nil {
		log.Fatal("File server not found.")
	} else {
		log.Printf("File server found.")
	}

	r := mux.NewRouter()

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/blindtest", games.HandleBlindtest)
	http.HandleFunc("/blindtestws", games.BlindtestWs)
	http.HandleFunc("/deaftest", games.HandleDeaftest)
	http.HandleFunc("/deaftestws", games.DeaftestWs)
	http.HandleFunc("/scattegories", games.HandleScattegories)
	http.HandleFunc("/ScattegoriesGame", games.HandleScattegoriesGameSocket)
	http.HandleFunc("/register", user.HandleRegister)
	http.HandleFunc("/login", user.HandleLogin)
	http.HandleFunc("/loginControl", user.HandleLoginControl)
	http.HandleFunc("/registerControl", user.HandleRegisterControl)
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/checkUsername", user.HandleCheckUsername) // Nouvelle route : checkUsername (pour vérifier la disponibilité d'un nom d'utilisateur lors de l'inscription par requête AJAX)
	http.HandleFunc("/checkEmail", user.HandleCheckEmail)       // Nouvelle route : checkEmail (pour vérifier la disponibilité d'un email lors de l'inscription par requête AJAX)
	http.HandleFunc("/isPasswordValid", user.HandleIsPasswordValid)
	http.HandleFunc("/lobby", lobby.HandleLobby)
	// http.HandleFunc("/logout", session.HandleLogout)
	http.HandleFunc("/getLetter", games.HandleGetLetter)
	http.HandleFunc("/createRoom", room.HandleCreateRoom)
	// http.HandleFunc("/joinRoom", room.HandleJoinRoom)
	r.HandleFunc("/room/{roomID}", games.HandleRoom)

	fmt.Println(time.Now().String() + " Server is running on port " + PORT)

	// Horloge pour les sessions
	// go func() {
	// 	for {
	// 		time.Sleep(1 * time.Second)
	// 		for _, session := range session.ActiveSessions {
	// 			if time.Now().Unix()-session.InactiveSince > 360 {
	// 				fmt.Println("Session inactive depuis 6 minutes. Suppression de la session.")
	// 				fmt.Println("Session : ", session)
	// 				// TODO
	// 			} else {
	// 				// Incrémenter la variable inactiveSince de chaque session de 1 seconde
	// 				session.InactiveSince++
	// 			}
	// 		}
	// 	}
	// }()

	server := &http.Server{
		Addr:              *addr,
		Handler:           nil,
		ReadHeaderTimeout: 3 * time.Second,
	}

	// Démarrage du serveur
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Printf("Server listening at " + *addr)
	log.Fatal(server.ListenAndServe())

	database.Database()

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/home-page.html")
}
