package models

import (
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

func CreateLetterInDB(letter Letter) error {
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

func GetLettersFromDB(username string) ([]Letter, error) {
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
