package games

import (
	spotifyapi "groupietracker/spotifyApi"
	mrand "math/rand"
	"net/http"
)

type PageData struct {
	Artist string
	Title  string
	URL    string
	Lyrics string
}

func HandleBlindtest(w http.ResponseWriter, r *http.Request) {
	var music spotifyapi.Music

	var artists = music.Artists
	var titles = music.Titles
	var musicUrl = music.MusicUrl

	i := mrand.Intn(len(musicUrl))

	data := PageData{
		Artist: artists[i],
		Title:  titles[i],
		URL:    musicUrl[i],
	}

	if r.Method == http.MethodPost {
		answer := r.FormValue("blindtest_answer")
		if answer == data.Title || answer == data.Artist || answer == data.Title+" "+data.Artist || answer == data.Artist+" "+data.Title {
			musicUrl = append(musicUrl[:i], musicUrl[i+1:]...)
			titles = append(titles[:i], titles[i+1:]...)
			artists = append(artists[:i], artists[i+1:]...)
			i = mrand.Intn(len(musicUrl))
			data.URL = musicUrl[i]
			data.Title = titles[i]
			data.Artist = artists[i]
		}

	}
	http.ServeFile(w, r, "blindtest.html")
}
