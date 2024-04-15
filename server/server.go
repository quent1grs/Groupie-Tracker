package main

import (
	"groupietracker/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.Database()
	database.GetAlbum("https://api.spotify.com/v1/playlists/3hhUZQwNteEDClZTu4XY9X")
}

// func getAccessToken() {
// 	clientID := ""
// 	clientSecret := ""

// 	// Base64 encode the client ID and secret
// 	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

// 	data := url.Values{}
// 	data.Set("grant_type", "client_credentials")

// 	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
// 	req.Header.Add("Authorization", "Basic "+auth)
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

// 	resp, _ := http.DefaultClient.Do(req)
// 	body, _ := ioutil.ReadAll(resp.Body)

// 	var result map[string]interface{}
// 	json.Unmarshal(body, &result)

// 	fmt.Println(result["access_token"])
// }
