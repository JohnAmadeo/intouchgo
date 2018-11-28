package models

import (
	"errors"
	"log"

	"github.com/johnamadeo/intouchgo/lob"
)

type Inmate struct {
	Id           string `json:"id"`
	State        string `json:"state"`
	InmateNumber string `json:"inmateNumber"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	DateOfBirth  string `json:"dateOfBirth"`
	Facility     string `json:"facility"`
	Active       bool   `json:"active"`
}

type InmateKey struct {
	State, InmateNumber string
}

type Facility struct {
	Name             string
	ShortName        string
	AddressLine      string
	City             string
	State            string
	Zip              string
	LobTestAddressId string
	LobLiveAddressId string
}

func getKey(inmate Inmate) InmateKey {
	return InmateKey{
		State:        inmate.State,
		InmateNumber: inmate.InmateNumber,
	}
}

func getInmatesKeySet(inmates []Inmate) map[InmateKey]interface{} {
	ids := make(map[InmateKey]interface{})

	for _, inmate := range inmates {
		ids[getKey(inmate)] = nil
	}

	return ids
}

func GetFacilitiesFromDB() ([]Facility, error) {
	facilities := []Facility{}

	db, err := getDBConnection()
	if err != nil {
		return facilities, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM facilities")
	if err != nil {
		return facilities, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, shortName, addressLine, city, state, zip, lobTestAddressId, lobLiveAddressId string
		err := rows.Scan(&name, &shortName, &addressLine, &city, &state, &zip, &lobTestAddressId, &lobLiveAddressId)

		if err != nil {
			return facilities, err
		}

		facilities = append(facilities, Facility{
			Name:             name,
			ShortName:        shortName,
			AddressLine:      addressLine,
			City:             city,
			State:            state,
			Zip:              zip,
			LobTestAddressId: lobTestAddressId,
			LobLiveAddressId: lobLiveAddressId,
		})
	}

	return facilities, nil
}

func GetInmateAddress(inmateId string, lobEnvironment string) (string, error) {
	db, err := getDBConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	addressType := ""
	switch lobEnvironment {
	case lob.LobTestEnvironment:
		addressType = "facilities.lobTestAddressId"
	case lob.LobLiveEnvironment:
		addressType = "facilities.lobLiveAddressId"
	default:
		return "", errors.New(lob.LobInvalidEnvironmentError)
	}

	rows, err := db.Query(
		"SELECT "+addressType+" FROM inmates JOIN facilities ON inmates.facility = facilities.name WHERE inmates.id = $1",
		inmateId,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	addressId := ""
	for rows.Next() {
		err := rows.Scan(&addressId)
		if err != nil {
			return "", err
		}

		break
	}

	if addressId == "" {
		return "", errors.New("No inmate with requested id found.")
	}
	return addressId, nil
}

func GetInmatesFromDB(searchQuery string) ([]Inmate, error) {
	inmates := []Inmate{}

	db, err := getDBConnection()
	if err != nil {
		return inmates, err
	}
	defer db.Close()

	query := "SELECT * " +
		"FROM inmates " +
		"WHERE UPPER(CONCAT(firstName, ' ', lastName)) LIKE UPPER('%' || $1 || '%')"

	rows, err := db.Query(query, searchQuery)
	if err != nil {
		return inmates, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, state, inmateNumber, firstName, lastName, dateOfBirth, facility string
		var active bool
		err := rows.Scan(
			&id,
			&state,
			&inmateNumber,
			&firstName,
			&lastName,
			&dateOfBirth,
			&facility,
			&active,
		)

		if err != nil {
			return inmates, err
		}

		inmate := Inmate{
			Id:           id,
			State:        state,
			InmateNumber: inmateNumber,
			FirstName:    firstName,
			LastName:     lastName,
			DateOfBirth:  dateOfBirth,
			Facility:     facility,
			Active:       active,
		}

		inmates = append(inmates, inmate)
	}

	return inmates, nil
}

func SaveInmatesFromScraper(scraperInmates []Inmate) error {
	dbInmates, err := GetInmatesFromDB("")
	if err != nil {
		return err
	}

	scraperIds := getInmatesKeySet(scraperInmates)
	dbIds := getInmatesKeySet(dbInmates)

	db, err := getDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Loop through all existing inmates in DB and mark them as inactive if they
	// are no longer on the website
	for _, dbInmate := range dbInmates {
		if _, ok := scraperIds[getKey(dbInmate)]; !ok {
			_, err = tx.Exec(
				"UPDATE inmates SET active = false WHERE state = $1 AND inmateNumber = $2",
				dbInmate.State,
				dbInmate.InmateNumber,
			)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// TODO: Batch inserts and updates for better performance

	// Insert all new inmates, and update existing inmates by marking them as
	// active and updating their facility
	for _, scraperInmate := range scraperInmates {
		if _, ok := dbIds[getKey(scraperInmate)]; !ok {
			_, err := tx.Exec(
				"INSERT INTO inmates VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
				scraperInmate.Id,
				scraperInmate.State,
				scraperInmate.InmateNumber,
				scraperInmate.FirstName,
				scraperInmate.LastName,
				scraperInmate.DateOfBirth,
				scraperInmate.Facility,
				scraperInmate.Active,
			)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			_, err = tx.Exec(
				"UPDATE inmates "+
					"SET active = true, facility = $1 "+
					"WHERE state = $2 AND inmateNumber = $3",
				scraperInmate.Facility,
				scraperInmate.State,
				scraperInmate.InmateNumber,
			)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
