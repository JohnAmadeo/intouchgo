package main

import (
	"log"
	"net/http"
	"os"

	"github.com/johnamadeo/intouchgo/auth"
	"github.com/johnamadeo/intouchgo/routes"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/inmates", auth.GetAuthHandler(routes.InmatesHandler))
	serveMux.Handle("/letter", auth.GetAuthHandler(routes.CreateLetterHandler))
	serveMux.Handle("/letters", auth.GetAuthHandler(routes.LettersHandler))
	serveMux.Handle("/user", auth.GetAuthHandler(routes.CreateUserHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, serveMux))
}
