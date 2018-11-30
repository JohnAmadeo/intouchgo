package models

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

const (
	Driver            = "postgres"
	LocalDBConnection = "user=johnamadeodaniswara dbname=intouch_android sslmode=disable"
)

func getDBConnection() (*sql.DB, error) {
	var dataSource string

	if dbURL, ok := os.LookupEnv("DATABASE_URL"); !ok {
		dataSource = LocalDBConnection
	} else {
		// the choice to use the "DATABASE_URL" environment variable is detailed in
		// https://devcenter.heroku.com/articles/heroku-postgresql#provisioning-heroku-postgres
		dataSource = dbURL
	}

	db, err := sql.Open(Driver, dataSource)
	if err != nil {
		db.Close()
		return db, err
	}

	return db, nil
}
