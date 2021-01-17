package postgres

import (
	"database/sql"
	"fmt"

	"../../common"
	"../../secrets"

	// This package is blank because although it is not
	// directly called, it is used as a driver
	_ "github.com/lib/pq"
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

// Insert inserts artists data into the database
func Insert(items []common.Item) {
	insertStatement := `INSERT INTO artists (id, name, popularity) VALUES ($1, $2, $3)`
	numInserted := 0
	for i := 0; i < len(items); i++ {
		_, err := db.Exec(insertStatement, items[i].ID, items[i].Name, items[i].Popularity)
		if err != nil {
			panic(err)
		}
		numInserted++
	}
	fmt.Println(fmt.Sprintf("%d row(s) inserted", numInserted))
}
