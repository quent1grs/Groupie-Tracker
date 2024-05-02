package games

import "net/http"

func HandleScattegories(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/scattegories.html")
}
