package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)
func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/letter", getAuthHandler(createLetterHandler))
	serveMux.Handle("/letters", getAuthHandler(lettersHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, serveMux))
}
