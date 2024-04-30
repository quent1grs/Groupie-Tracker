package games

import "net/http"

func HandleScattegories(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./scattegories.html")
}
