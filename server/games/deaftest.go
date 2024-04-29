package games

import (
	spotifyapi "groupietracker/spotifyApi"
	mrand "math/rand"
	"net/http"
)

func HandleDeaftest(w http.ResponseWriter, r *http.Request) {
	var music spotifyapi.Music

	var artists = music.Artists
	var titles = music.Titles
	var musicLyrics = music.MusicLyrics

	i := mrand.Intn(len(musicLyrics))

	data := PageData{
		Artist: artists[i],
		Title:  titles[i],
		Lyrics: musicLyrics[i],
	}

	if r.Method == http.MethodPost {
		answer := r.FormValue("deaftest_answer")
		if answer == data.Title || answer == data.Artist || answer == data.Title+" "+data.Artist || answer == data.Artist+" "+data.Title {
			musicLyrics = append(musicLyrics[:i], musicLyrics[i+1:]...)
			titles = append(titles[:i], titles[i+1:]...)
			artists = append(artists[:i], artists[i+1:]...)
			i = mrand.Intn(len(musicLyrics))
			data.Lyrics = musicLyrics[i]
			data.Title = titles[i]
			data.Artist = artists[i]
		}
	}

	http.ServeFile(w, r, "./deaftest.html")
}
