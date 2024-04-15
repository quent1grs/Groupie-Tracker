package main

import (
	"encoding/base64"
	"encoding/json"
	"groupietracker/database"
	"io"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	token := getToken()
	database.Database()
	database.GetAlbum("https://api.spotify.com/v1/playlists/3hhUZQwNteEDClZTu4XY9X", token)

}

func getToken() string {
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
