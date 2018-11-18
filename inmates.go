package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Inmate struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	InmateNumber string `json:"inmateNumber"`
	DateOfBirth  string `json:"dateOfBirth"`
	Facility     string `json:"facility"`
}

func getInmatesFromDB(searchQuery string) ([]Inmate, error) {
	inmates := []Inmate{}

	db, err := getDBConnection()
	if err != nil {
		return inmates, err
	}
	defer db.Close()

	fields := []string{
		"inmates.id",
		"CONCAT(inmates.firstName, ' ', inmates.lastName) AS name",
		"inmates.inmateNumber",
		"inmates.dateOfBirth",
		"inmates.facility",
	}

	query := "SELECT * FROM (" +
		"SELECT " + strings.Join(fields[:], ", ") + " " +
		"FROM inmates " +
		") AS inmates " +
		"WHERE UPPER(name) LIKE UPPER('%' || $1 || '%')"

	rows, err := db.Query(query, searchQuery)
	if err != nil {
		return inmates, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, inmateNumber, dateOfBirth, facility string
		err := rows.Scan(
			&id,
			&name,
			&inmateNumber,
			&dateOfBirth,
			&facility,
		)

		if err != nil {
			return inmates, err
		}

		inmate := Inmate{
			Id:           id,
			Name:         name,
			InmateNumber: inmateNumber,
			DateOfBirth:  dateOfBirth,
			Facility:     facility,
		}

		inmates = append(inmates, inmate)
	}

	return inmates, nil
}

func inmatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(messageToBytes("Only GET requests are allowed at this route"))
		return
	}

	queries, ok := r.URL.Query()["query"]
	if !ok || len(queries) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(messageToBytes("Request query parameters must contain a single username"))
		return
	}

	inmates, err := getInmatesFromDB(queries[0])
	if err != nil {
		printErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(messageToBytes(err.Error()))
		return
	}

	bytes, err := json.Marshal(inmates)
	if err != nil {
		printErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(messageToBytes(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
