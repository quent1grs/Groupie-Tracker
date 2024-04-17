package websocket

import (
	"fmt"
	"groupietracker/database"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Mise à niveau de la connexion HTTP en une connexion WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erreur lors de l'upgrade de la connexion WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		// Lecture des messages WebSocket
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erreur lors de la lecture du message WebSocket:", err)
			return
		}

		// Vérification du type de message (supposons que les données du formulaire sont envoyées en texte)
		if messageType == websocket.TextMessage {
			// Traitement des données du formulaire
			formData := strings.Split(string(p), ",") // Supposons que les données du formulaire sont séparées par des virgules
			email := formData[0]
			password := formData[1]
			username := formData[2]

			// Insertion des données dans la base de données
			err := database.InsertFormData(email, password, username)
			if err != nil {
				fmt.Println("Erreur lors de l'insertion des données dans la base de données:", err)
				// Gérer l'erreur de manière appropriée (par exemple, envoyer une réponse d'erreur au client)
				return
			}
		}
	}
}
