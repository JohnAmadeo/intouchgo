package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Message string
}

func messageToBytes(message string) []byte {
	bytes, _ := json.Marshal(Message{message})
	return bytes
}

func printErr(err error) {
	fmt.Println(err.Error())
}
