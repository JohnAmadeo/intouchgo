package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/johnamadeo/intouchgo/auth"
	"github.com/johnamadeo/intouchgo/models"
	"github.com/johnamadeo/intouchgo/routes"
	"github.com/johnamadeo/intouchgo/utils"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/inmates", auth.GetAuthHandler(routes.InmatesHandler))
	serveMux.Handle("/letter", auth.GetAuthHandler(routes.CreateLetterHandler))
	serveMux.Handle("/letters", auth.GetAuthHandler(routes.LettersHandler))
	serveMux.Handle("/user", auth.GetAuthHandler(routes.CreateUserHandler))
	serveMux.Handle("/", http.FileServer(http.Dir("./static")))

	serveMux.Handle("/test/facilities", auth.GetFakeAuthHandler(func(w http.ResponseWriter, r *http.Request) {
		facilities, err := models.GetFacilitiesFromDB()
		if err != nil {
			utils.PrintErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.MessageToBytes(err.Error()))
			return
		}

		bytes, err := json.Marshal(facilities)
		if err != nil {
			utils.PrintErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.MessageToBytes(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, serveMux))
}
