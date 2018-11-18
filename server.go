package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/inmates", getAuthHandler(inmatesHandler))
	serveMux.Handle("/letter", getAuthHandler(createLetterHandler))
	serveMux.Handle("/letters", getAuthHandler(lettersHandler))
	serveMux.Handle("/user", getAuthHandler(createUserHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, serveMux))
}
