package games

import "net/http"

func HandleDeafTest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./deafTest.html")
}
