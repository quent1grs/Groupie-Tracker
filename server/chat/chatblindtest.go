package chat

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

func ChatBlindtest() {
	hub := []*websocket.Conn{}

	http.Handle("/chatblindtestws", websocket.Handler(func(ws *websocket.Conn) {
		data := map[string]interface{}{}
		hub = append(hub, ws)
		for {
			err := websocket.JSON.Receive(ws, &data)
			if err != nil {
				fmt.Println("Error reading json.", err)
				ws.Close()
				break
			}
			fmt.Println("Received data: ", data)
			message := fmt.Sprintf("dit : %s", data["message"])
			for i := range hub {
				websocket.Message.Send(hub[i], message)
			}
		}

	}))
}
