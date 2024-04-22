package spotifyapi

import (
	"io"
	"log"
	"net/http"

	mm "github.com/fr05t1k/musixmatch"
	"github.com/fr05t1k/musixmatch/entity/lyrics"
	mmhttp "github.com/fr05t1k/musixmatch/http"
)

type Track struct {
	Album        Album        `json:"album"`
	Name         string       `json:"name"`
	Artists      []Artist     `json:"artists"`
	Genres       []Genres     `json:"genres"`
	ExternalUrls ExternalUrls `json:"external_urls"`
}

type Genres struct {
	Genres []Genres `json:"genres"`
}

type Item struct {
	Track Track  `json:"track"`
	Name  string `json:"name"`
}

type Artist struct {
	Name string `json:"name"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

type SearchResponse struct {
	Tracks struct {
		Items []Item `json:"items"`
	} `json:"tracks"`
}

type Album struct {
	Name string `json:"name"`
}

type LyricsBody struct {
	Lyrics []Lyric `json:"body"`
}

type Lyric struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Language string `json:"language"`
}

func GetPlaylist(url string, token string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func GetLyrics(title string, artist string) *lyrics.Lyrics {
	client := mm.NewClient("78b38fd30f412e2735ef3229e3f93e94")
	println("title : " + title + " artist : " + artist)

	searchRequest := mmhttp.SearchRequest{QueryTrack: title, QueryArtist: artist}

	tracks, err := client.SearchTrack(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	var lyrics *lyrics.Lyrics
	if len(tracks) > 0 {
		trackID := tracks[0].Track.Id

		lyrics, _ = client.GetLyrics(trackID)

	}
	return lyrics
}
