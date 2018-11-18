package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	Driver       = "postgres"
	DBConnection = "user=johnamadeodaniswara dbname=intouch_android sslmode=disable"
)

func getDBConnection() (*sql.DB, error) {
	db, err := sql.Open(Driver, DBConnection)

	if err != nil {
		db.Close()
		return db, err
	}

	return db, nil
}
