package utils

import (
	"encoding/json"
	"log"
)

func WriteJson(data interface{}) (bytes []byte) {
	log.Println("marshaling the data ", data)
	bytes, _ = json.Marshal(data)

	return bytes
}
