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

// Close closes the database in the shutdown hook
func Close() {
	db.Close()
}

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
	insertStatement := `INSERT INTO artists (id, name, popularity, image_data) VALUES ($1, $2, $3, $4)`
	numInserted := 0
	for i := 0; i < len(items); i++ {
		popularity := items[i].Popularity
		if popularity <= 60 {
			continue
		}
		id := items[i].ID
		var image common.Image
		if len(items[i].Images) > 0 {
			image = items[i].Images[0]
		}
		_, err := db.Exec(insertStatement, id, items[i].Name, items[i].Popularity, image)
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key") {
				panic(err)
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
		if recommendedArtists[i].Popularity <= 60 {
			continue
		}
		artists = append(artists, recommendedArtists[i].ID)
	}
	_, err := db.Exec(updateStatement, pq.Array(artists), id)
	if err != nil {
		panic(err)
	}
}

// UpdateMissingRecommended essentially makes an undirected graph
func UpdateMissingRecommended(m map[string]common.Array) {
	updateStatement := "UPDATE artists SET recommended = $1 WHERE id = $2"
	for k, v := range m {
		_, err := db.Exec(updateStatement, pq.Array(v), k)
		if err != nil {
			panic(err)
		}
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

// GetAllArtists returns all of the artists
func GetAllArtists() []common.Artist {
	rows, err := db.Query("SELECT id, name, popularity, recommended, image_data FROM artists")
	if err != nil {
		panic(err)
	}
	var artists []common.Artist
	for rows.Next() {
		var artist common.Artist
		err = rows.Scan(&artist.ID, &artist.Name, &artist.Popularity, pq.Array(&artist.Recommended), &artist.ImageData)
		artists = append(artists, artist)
	}
	return artists
}
