package models

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/johnamadeo/intouchgo/utils"
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

/*
curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlJqVTJRVEZFTTBVeU5EazFRalk1TmtFM05UazVOak0xUVVJeVFUWTRPVEZGUVVJeFJEY3lOdyJ9.eyJpc3MiOiJodHRwczovL2ludG91Y2gtYW5kcm9pZC5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NWJkODk1MWEzZGRjYjQwNWQzYWM2Y2RjIiwiYXVkIjpbImh0dHBzOi8vaW50b3VjaC1hbmRyb2lkLWJhY2tlbmQuaGVyb2t1YXBwLmNvbS8iLCJodHRwczovL2ludG91Y2gtYW5kcm9pZC5hdXRoMC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNTQyNTA3MTM2LCJleHAiOjE1NDI1OTM1MzYsImF6cCI6InpjVU54OGxQQVE2djhVSXB4OTIwVkdvVTVnMmplNXl6Iiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCBvZmZsaW5lX2FjY2VzcyJ9.kNFhzrMJ75jPjWufClOm8oXYTGIGEZHnmDGOV0Hu-T8r_ObtFhngsT1QoqI8-6Y2jIyKJyKUc4ySbfS9PRAeYaulntnFlvdcJXJK6M5_gNizSS84ROgb6vHVBxdZ6YfIIS3iQbf51g2xJs-jRWWwyb1_Shr8Fnt0fxOu_D_KBt82GiQWdJgePT-SVKMRrmZno_4aq07-YFxdJppLbfQIeMz4a4kApBHQ1ZKPZ6qrUcZ1xPKqd3dTMXrMtWAvkISGp31zTPEbm-40t4xIWgfHOXK20WJ-KZ4zV7nczmYHMNqdxx2ww5J9xc2vXah2ITPJagQUlQH6ZqrOsmGC4ZJcGA" http://localhost:8080/inmates
*/
func InmatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(utils.MessageToBytes("Only GET requests are allowed at this route"))
		return
	}

	queries, ok := r.URL.Query()["query"]
	if !ok || len(queries) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.MessageToBytes("Request query parameters must contain a single username"))
		return
	}

	inmates, err := getInmatesFromDB(queries[0])
	if err != nil {
		utils.PrintErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MessageToBytes(err.Error()))
		return
	}

	bytes, err := json.Marshal(inmates)
	if err != nil {
		utils.PrintErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MessageToBytes(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
