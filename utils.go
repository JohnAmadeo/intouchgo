package main

import (
	"encoding/json"
)

type Message struct {
	Message string
}

func messageToBytes(message string) []byte {
	bytes, _ := json.Marshal(Message{message})
	return bytes
}
