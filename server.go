package main

import (
	"log"
	"net/http"
	"os"

	"github.com/johnamadeo/intouchgo/auth"
	"github.com/johnamadeo/intouchgo/models"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/inmates", auth.GetAuthHandler(models.InmatesHandler))
	serveMux.Handle("/letter", auth.GetAuthHandler(models.CreateLetterHandler))
	serveMux.Handle("/letters", auth.GetAuthHandler(models.LettersHandler))
	serveMux.Handle("/user", auth.GetAuthHandler(models.CreateUserHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, serveMux))
}
