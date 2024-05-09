package games

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

type Letter struct {
	Letter string `json:"letter"`
}

func HandleScattegories(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/scattegories.html")
}

func HandleScattegoriesGameSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var gameAction map[string]string
		err = json.Unmarshal(message, &gameAction)
		if err != nil {
			log.Println("json unmarshal:", err)
			break
		}

		if gameAction["action"] == "stopGame" {
			// Créez un message pour dire aux clients de stopper le jeu
			stopMessage := map[string]string{
				"action": "stopGame",
			}
			stopMessageJson, err := json.Marshal(stopMessage)
			if err != nil {
				log.Println("json marshal:", err)
				break
			}

			//FAIRE EN FONCTION DES ROOMS USER
			// Remplacez `clients` par la liste de vos clients WebSocket
			for _, client := range clients {
				if err := client.WriteMessage(websocket.TextMessage, stopMessageJson); err != nil {
					log.Println("write:", err)
				}
			}
		}
	}
}

func HandleGetLetter(w http.ResponseWriter, r *http.Request) {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letter := string(letters[rand.Intn(len(letters))])
	json.NewEncoder(w).Encode(Letter{Letter: letter})
}
