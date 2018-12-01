package models

import (
	"strings"

	"github.com/johnamadeo/intouchgo/lob"
	"github.com/johnamadeo/intouchgo/utils"
)

type Letter struct {
	Id                    string `json:"id"`
	Author                string `json:"author"`
	Recipient             string `json:"recipient"`
	RecipientId           string `json:"recipientId"`
	Subject               string `json:"subject"`
	Text                  string `json:"text"`
	TimeSent              string `json:"timeSent"`
	TimeLastEdited        string `json:"timeLastEdited"`
	TimeDeliveredEstimate string `json:"timeDeliveredEstimate"`
	IsDraft               bool   `json:"isDraft"`
	LobLetterId           string `json:"lobLetterId"`
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
		"TO_CHAR(letters.timeDeliveredEstimate, 'MM/dd/yy')",
		"letters.isDraft",
		"letters.lobLetterId",
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
		var id, author, recipient, recipientId, subject, text, timeSent, timeLastEdited, timeDeliveredEstimate, lobLetterId string
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
			&timeDeliveredEstimate,
			&isDraft,
			&lobLetterId,
		)

		if err != nil {
			return letters, err
		}

		letter := Letter{
			Id:                    id,
			Author:                author,
			Recipient:             recipient,
			RecipientId:           recipientId,
			Subject:               subject,
			Text:                  text,
			TimeSent:              timeSent,
			TimeLastEdited:        timeLastEdited,
			TimeDeliveredEstimate: timeDeliveredEstimate,
			IsDraft:               isDraft,
			LobLetterId:           lobLetterId,
		}

		letters = append(letters, letter)
	}

	return letters, nil
}

func SendLetter(letter Letter) (Letter, error) {
	response, err := sendLetterToLob(letter, lob.LobTestEnvironment)
	if err != nil {
		return Letter{}, err
	}

	date, err := lob.LobDateToDBDate(response.ExpectedDeliveryDate)
	if err != nil {
		return Letter{}, err
	}

	letter.TimeDeliveredEstimate = date
	letter.LobLetterId = response.Id

	err = createLetterInDB(letter)
	if err != nil {
		return Letter{}, err
	}

	return letter, nil
}

func createLetterInDB(letter Letter) error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO letters VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		letter.Id,
		letter.Author,
		letter.RecipientId,
		letter.Subject,
		letter.Text,
		letter.TimeSent,
		letter.TimeLastEdited,
		letter.TimeDeliveredEstimate,
		letter.IsDraft,
		letter.LobLetterId,
	)

	if err != nil {
		return err
	}

	return nil
}

func sendLetterToLob(letter Letter, lobEnvironment string) (lob.LobSendLetterResponse, error) {
	var response lob.LobSendLetterResponse

	inmateAddressId, err := GetInmateAddress(letter.RecipientId, lobEnvironment)
	if err != nil {
		return response, err
	}

	inTouchAddressId, err := lob.GetInTouchAddress(lobEnvironment)
	if err != nil {
		return response, err
	}

	htmlString, err := lob.GetLetterHTMLTemplate(letter.Text)
	if err != nil {
		return response, err
	}

	request := lob.LobSendLetterRequest{
		Color:    false,
		MailType: lob.USPSStandard,
		From:     inTouchAddressId,
		To:       inmateAddressId,
		File:     htmlString,
		MergeVariables: map[string]string{
			"author":    letter.Author,
			"recipient": letter.Recipient,
			"subject":   letter.Subject,
			"timeSent":  letter.TimeSent,
		},
	}

	err = lob.Post(
		"letters",
		utils.JSONToForm(request),
		&response,
		lob.LobTestEnvironment,
	)
	if err != nil {
		return response, err
	}

	return response, nil
}
