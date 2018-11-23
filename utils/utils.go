package utils

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Message string
}

func MessageToBytes(message string) []byte {
	bytes, _ := json.Marshal(Message{message})
	return bytes
}

func PrintErr(err error) {
	fmt.Println(err.Error())
}
