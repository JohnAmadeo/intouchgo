package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/johnamadeo/intouchgo/models"
	"github.com/johnamadeo/intouchgo/utils"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(utils.MessageToBytes("Only POST requests are allowed at this route"))
		return
	}

	var user models.User
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.MessageToBytes("Malformed body."))
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(bytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.MessageToBytes("Request body must be a user."))
		return
	}

	err = models.CreateUser(user)
	if err != nil {
		utils.PrintErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MessageToBytes("Failed to create user: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(utils.MessageToBytes("Successfully created user."))
}
