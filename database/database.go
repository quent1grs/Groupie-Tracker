package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Database() {
	// user_id := 0
	db, err := sql.Open("sqlite3", "/Users/quentingros/Desktop/Groupie-Tracker/database/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	GetUsers(db)
	isEmpty, err := isUserTableEmpty(db)
	if err != nil {
		log.Fatal(err)
	}

	if isEmpty {
		// user_id = 0
	} else {
		// user_id, err = getNextID(db)
		if err != nil {
			log.Fatal(err)
		}

	}
}

// 	fmt.Println("do you want to register ? (y/n)")
// 	var response string
// 	fmt.Scanln(&response)

// 	if response == "y" {
// 		register(db, user_id)
// 	}

// 	fmt.Println("do you want to delete all users ? (y/n)")
// 	fmt.Scanln(&response)

// 	if response == "y" {
// 		deleteAllUsers(db)
// 	}

// 	rows, err := db.Query("SELECT * FROM USER")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var id int
// 		var pseudo string
// 		var email string
// 		var password string
// 		err = rows.Scan(&id, &pseudo, &email, &password)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(id, pseudo, email, password)
// 	}
// }

// InsertFormData insère les données du formulaire dans la base de données
func InsertFormData(username, email, password string) error {
	// Ouverture de la connexion à la base de données
	db, err := sql.Open("sqlite3", "./database/db.sqlite")
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture de la base de données: %v", err)
	}

	// Préparation de la requête d'insertion
	stmt, err := db.Prepare("INSERT INTO USER(username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("erreur lors de la préparation de la requête d'insertion: %v", err)
	}
	defer stmt.Close()

	// Exécution de la requête d'insertion avec les valeurs du formulaire
	_, err = stmt.Exec(username, email, password)
	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion dans la base de données: %v", err)
	}

	fmt.Println("Données du formulaire insérées avec succès dans la base de données.")
	db.Close()
	return nil

}

// func register(db *sql.DB, id int) {
// 	if id > 1 {
// 		id = -1
// 	}
// 	pseudo := ""
// 	email := ""
// 	password := ""

// 	fmt.Println("Enter your pseudo: ")
// 	fmt.Scanln(&pseudo)
// 	fmt.Println("Enter your email: ")
// 	fmt.Scanln(&email)
// 	fmt.Println("Enter your password: ")
// 	fmt.Scanln(&password)

// 	id++
// }

func GetUsers(db *sql.DB) {
	// Préparation de la requête SQL
	rows, err := db.Query("SELECT id, username, email, password FROM USER")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Parcours des lignes de résultat
	for rows.Next() {
		var id int
		var username, email, password string

		// Scan des colonnes dans les variables
		err = rows.Scan(&id, &username, &email, &password)
		if err != nil {
			log.Fatal(err)
		}

		// Affichage des données
		fmt.Println(id, username, email, password)
	}

	// Vérification d'erreurs après avoir parcouru les lignes
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID(db *sql.DB) (int, error) {
	var maxID int
	err := db.QueryRow("SELECT MAX(id) FROM USER").Scan(&maxID)
	if err != nil {
		return 0, err
	}
	return maxID + 1, nil
}

func isUserTableEmpty(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM USER").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func deleteAllUsers(db *sql.DB) {
	stmt, err := db.Prepare("DELETE FROM USER")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

type Track struct {
	Name     string   `json:"name"`
	Artists  []Artist `json:"artists"`
	Genre    []Genre  `json:"genre"`
	Duration int      `json:"duration_ms"`
}

type Item struct {
	Track Track `json:"track"`
}

type Genre struct {
	Genre []Genre `json:"genre"`
}

type Artist struct {
	Name string `json:"name"`
}

type SearchResponse struct {
	Tracks struct {
		Items []Item `json:"items"`
	} `json:"tracks"`
}

type duration_ms struct {
	Duration int `json:"duration_ms"`
}

func GetPlaylist(url string, token string) {
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
	printTitle(body, err)
}

func printTitle(body []byte, err error) {
	fmt.Println("Title		|	Artists")
	var searchResponse SearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range searchResponse.Tracks.Items {
		fmt.Print(item.Track.Name + "	| ")
		for _, artist := range item.Track.Artists {
			fmt.Print(artist.Name, ", ")
		}
		fmt.Println()
	}
}
