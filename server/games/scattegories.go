package games

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type Letter struct {
	Letter string `json:"letter"`
}

func HandleScattegories(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/scattegories.html")
}

func HandleGetLetter(w http.ResponseWriter, r *http.Request) {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letter := string(letters[rand.Intn(len(letters))])
	json.NewEncoder(w).Encode(Letter{Letter: letter})
}
