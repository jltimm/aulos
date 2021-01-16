package postgres

import (
	"database/sql"
	"fmt"

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

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

}
