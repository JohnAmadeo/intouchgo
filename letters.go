package main

import (
	"encoding/json"
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

func getLettersFromDB() ([]Letter, error) {
	var letters []Letter

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
		"ON letters.recipient = inmates.id"

	rows, err := db.Query(query)
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

func lettersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(messageToBytes("Only GET requests are allowed at this route"))
		return
	}

	letters, err := getLettersFromDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(messageToBytes(err.Error()))
		return
	}

	bytes, err := json.Marshal(letters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(messageToBytes(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
