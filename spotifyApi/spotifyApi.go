package spotifyapi

import (
	"io"
	"log"
	"net/http"
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
	//printTitle(body)
	return body
}
