package spotifyapi

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	mm "github.com/fr05t1k/musixmatch"
	"github.com/fr05t1k/musixmatch/entity/lyrics"
	mmhttp "github.com/fr05t1k/musixmatch/http"
)

type Music struct {
	Artists     []string
	Titles      []string
	MusicUrl    []string
	MusicLyrics []string
}

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

func GetPlaylist(url string) []byte {
	client := &http.Client{}
	token := GetToken()

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

func GetToken() string {
	clientID := "c27ae1942ee94d23a21f324b6feba015"
	clientSecret := "c527485ba55545a4a0e88614a886500a" // Base64 encode the client ID and secret
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return result["access_token"].(string)
}

func ParsePlaylist(body []byte) {
	var music Music
	var playlist SearchResponse
	err := json.Unmarshal(body, &playlist)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range playlist.Tracks.Items {
		uri := path.Base(item.Track.ExternalUrls.Spotify)
		music.MusicUrl = append(music.MusicUrl, uri)
		title := item.Track.Name
		music.Titles = append(music.Titles, title)
		artist := item.Track.Artists[0].Name
		music.Artists = append(music.Artists, artist)
		lyrics := GetLyrics(title, artist)
		parts := strings.Split(lyrics.Body, "\n...\n\n******* This Lyrics is NOT for Commercial use *******")
		lyrics.Body = parts[0]
		music.MusicLyrics = append(music.MusicLyrics, lyrics.Body)
		println("Title: " + title + " Artist: " + artist)
	}
}
