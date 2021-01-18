package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"../../common"
	"../../secrets"

	"github.com/lib/pq"
)

var db *sql.DB

// Initialize postgres
func Initialize() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		secrets.GetHost(), secrets.GetPort(), secrets.GetUser(), secrets.GetPassword(), secrets.GetDbname())
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

// InsertArtists inserts artists data into the database
func InsertArtists(items []common.Item) {
	insertStatement := `INSERT INTO artists (id, name, popularity) VALUES ($1, $2, $3)`
	numInserted := 0
	for i := 0; i < len(items); i++ {
		id := items[i].ID
		_, err := db.Exec(insertStatement, id, items[i].Name, items[i].Popularity)
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key") {
				panic(err)
			} else {
				log.Printf("Duplicate key found: %s\n", id)
			}
		} else {
			numInserted++
		}
	}
	log.Printf("%d row(s) inserted\n", numInserted)
}

// UpdateRecommended updates an artists row with their
// recommended artists
func UpdateRecommended(id string, recommendedArtists []common.Item) {
	updateStatement := "UPDATE artists SET recommended = $1 WHERE id = $2"
	var artists []string
	for i := 0; i < len(recommendedArtists); i++ {
		artists = append(artists, recommendedArtists[i].ID)
	}
	_, err := db.Exec(updateStatement, pq.Array(artists), id)
	if err != nil {
		panic(err)
	}
}

// GetArtistCount returns the number of rows in the artists table
func GetArtistCount() int {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM artists")
	err := row.Scan(&count)
	if err != nil {
		panic(err)
	}
	return count
}

// GetArtistsWithNullRecommended returns artists where the recommended column is blank
func GetArtistsWithNullRecommended() []string {
	rows, err := db.Query("SELECT id FROM artists WHERE recommended IS NULL")
	if err != nil {
		panic(err)
	}
	var artists []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		artists = append(artists, id)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return artists
}
