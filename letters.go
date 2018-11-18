package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Letter struct {
	Id             string `json:"id"`
	Author         string `json:"author"`
	Recipient      string `json:"recipient"`
	RecipientId    string `json:"recipientId"`
	Subject        string `json:"subject"`
	Text           string `json:"text"`
	TimeSent       string `json:"timeSent"`
	TimeLastEdited string `json:"timeLastEdited"`
	IsDraft        bool   `json:"isDraft"`
}

func createLetterInDB(letter Letter) error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO letters VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		letter.Id,
		letter.Author,
		letter.RecipientId,
		letter.Subject,
		letter.Text,
		letter.TimeSent,
		letter.TimeLastEdited,
		letter.IsDraft,
	)

	if err != nil {
		return err
	}

	return nil
}

func getLettersFromDB(username string) ([]Letter, error) {
	letters := []Letter{}

	db, err := getDBConnection()
	if err != nil {
		return letters, err
	}
	defer db.Close()

	fields := []string{
		"letters.id",
		"letters.author",
		"CONCAT(inmates.firstName, ' ', inmates.lastName) AS recipient",
		"letters.recipient AS recipientId",
		"letters.subject",
		"letters.text",
		"TO_CHAR(letters.timeSent, 'MM/dd/yy')",
		"TO_CHAR(letters.timeLastEdited, 'MM/dd/yy')",
		"letters.isDraft",
	}

	query := "SELECT " + strings.Join(fields[:], ", ") + " " +
		"FROM letters JOIN inmates " +
		"ON letters.recipient = inmates.id " +
		"WHERE letters.author = $1"

	rows, err := db.Query(query, username)
	if err != nil {
		return letters, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, author, recipient, recipientId, subject, text, timeSent, timeLastEdited string
		var isDraft bool
		err := rows.Scan(
			&id,
			&author,
			&recipient,
			&recipientId,
			&subject,
			&text,
			&timeSent,
			&timeLastEdited,
			&isDraft,
		)

		if err != nil {
			return letters, err
		}

		letter := Letter{
			Id:             id,
			Author:         author,
			Recipient:      recipient,
			RecipientId:    recipientId,
			Subject:        subject,
			Text:           text,
			TimeSent:       timeSent,
			TimeLastEdited: timeLastEdited,
			IsDraft:        isDraft,
		}

		letters = append(letters, letter)
	}

	return letters, nil
}

func createLetterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(messageToBytes("Only POST requests are allowed at this route"))
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(messageToBytes("Malformed body."))
		return
	}
	defer r.Body.Close()

	var letter Letter
	err = json.Unmarshal(bytes, &letter)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(messageToBytes("Request body must be a letter"))
		return
	}

	err = createLetterInDB(letter)
	if err != nil {
		printErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(messageToBytes("Error inserting to DB: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(messageToBytes("Letter successfully created!"))
}

func lettersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(messageToBytes("Only GET requests are allowed at this route"))
		return
	}

	usernames, ok := r.URL.Query()["username"]
	if !ok || len(usernames) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(messageToBytes("Request query parameters must contain a single username"))
		return
	}

	letters, err := getLettersFromDB(usernames[0])
	if err != nil {
		printErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(messageToBytes(err.Error()))
		return
	}

	bytes, err := json.Marshal(letters)
	if err != nil {
		printErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(messageToBytes(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
